package dstash

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"
    "github.com/golang/protobuf/ptypes/wrappers"
    pb "github.com/guymolinari/dstash/grpc"
)

type StringSearch struct {
	*Conn
    client				 []pb.StringSearchClient
    indexBatch           map[string]struct{}
	batchSize			 int
	batchMutex			 sync.Mutex
}


func NewStringSearch(conn *Conn, batchSize int) *StringSearch {

    clients := make([]pb.StringSearchClient, len(conn.clientConn))
    for i := 0; i < len(conn.clientConn); i++ {
        clients[i] = pb.NewStringSearchClient(conn.clientConn[i])
    }
    return &StringSearch{Conn: conn, batchSize: batchSize, client: clients}
}


func (c *StringSearch) Dispose() error {

    c.batchMutex.Lock()
    defer c.batchMutex.Unlock()

    if c.indexBatch != nil {
        if err := c.BatchIndex(c.indexBatch); err != nil {
            return err
        }   
        c.indexBatch = nil
    } 
	return nil
}


func (c *StringSearch) splitStringBatch(batch map[string]struct{}, replicas int) []map[string]struct{} {

	c.Conn.nodeMapLock.RLock()
	defer c.Conn.nodeMapLock.RUnlock()

    batches := make([]map[string]struct{}, len(c.client))
    for i, _ := range batches {
        batches[i] = make(map[string]struct{}, 0)
    }
    for k, v := range batch {
        nodeKeys := c.Conn.hashTable.GetN(replicas, ToString(k))
    	// Iterate over node key list and collate into batches
    	for _, nodeKey := range nodeKeys  {
			if i, ok := c.Conn.nodeMap[nodeKey]; ok {
				batches[i][k] = v
			}
		}
	}
	return batches
}


func (c *StringSearch) Index(str string) error {
  
    c.batchMutex.Lock()
    defer c.batchMutex.Unlock()

    if c.indexBatch == nil {
		c.indexBatch = make(map[string]struct{}, 0)
	}

    c.indexBatch[str] = struct{}{}

    if len(c.indexBatch) >= c.batchSize {
		if err := c.BatchIndex(c.indexBatch); err != nil {
			return err
		}
		c.indexBatch = nil
    }
	return nil
}
   

func (c *StringSearch) BatchIndex(batch map[string]struct{}) error {

    batches := c.splitStringBatch(batch, c.Conn.Replicas)

    done := make(chan error)
    defer close(done)
    count := len(batches)
    for i, v := range batches {
        go func(client pb.StringSearchClient, m map[string]struct{}) {
			done <- c.batchIndex(client, m)
		}(c.client[i], v)
	}
	for {
		err := <- done
		if err != nil {
			return err
		}
		count--
		if count == 0 {
			break
		}
	}
	return nil
}


func (c *StringSearch) batchIndex(client pb.StringSearchClient, batch map[string]struct{}) error {

    ctx, cancel := context.WithTimeout(context.Background(), Deadline * time.Second)
    defer cancel()
    b := make([]*wrappers.StringValue, len(batch))
    i := 0
    for k, _ := range batch {
	    b[i] = &wrappers.StringValue{ Value: k }
		i++
    }
    stream, err := client.BatchIndex(ctx)
    if err != nil {
        return fmt.Errorf("%v.BatchIndex(_) = _, %v: ", c.client, err)
    }

	for i := 0; i < len(b); i++ {
        if err := stream.Send(b[i]); err != nil {
            return fmt.Errorf("%v.Send(%v) = %v", stream, b[i], err)
        }
	}
    _, err2 := stream.CloseAndRecv()
    if err2 != nil {
        return fmt.Errorf("%v.CloseAndRecv() got error %v, want %v", stream, err2, nil)
    }
    return nil
}


func (c *StringSearch) Search(searchTerms string) (map[uint64]struct{}, error) {

    results := make(map[uint64]struct{}, 0)
    rchan := make(chan map[uint64]struct{})
    defer close(rchan)
    done := make(chan error)
    defer close(done)
    count := len(c.client)

    for i, _ := range c.client {
        go func(client pb.StringSearchClient) {
            r, e := c.search(client, searchTerms)
            rchan <- r
            done <- e
        }(c.client[i])
    }

    for { 
        b := <- rchan
        if b != nil {
            for k, v := range b {
                results[k] = v
            }
        }
        err := <- done
        if err != nil {
            return nil, err
        }
        count--
        if count == 0 {
            break
        }
    }

    return results, nil
}


func (c *StringSearch) search(client pb.StringSearchClient, searchTerms string) (map[uint64]struct{}, error) {

    batch := make(map[uint64]struct{}, 0)

    ctx, cancel := context.WithTimeout(context.Background(), Deadline)
    defer cancel()

    stream, err := client.Search(ctx, &wrappers.StringValue{Value: searchTerms})
    if err != nil {
        return nil, fmt.Errorf("%v.Search(_) = _, %v", client, err)
    }
    for {
        item, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("%v.Search(_) = _, %v", client, err)
        }
		batch[item.Value] = struct{}{}
    }

    return batch, nil
}


