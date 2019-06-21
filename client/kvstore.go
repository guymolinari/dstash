package dstash

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"time"
    "github.com/golang/protobuf/ptypes/empty"
    pb "github.com/guymolinari/dstash/grpc"
)

type KVStore struct {
	*Conn
    client				 []pb.KVStoreClient
}


func NewKVStore(conn *Conn) *KVStore {

    clients := make([]pb.KVStoreClient, len(conn.clientConn))
    for i := 0; i < len(conn.clientConn); i++ {
        clients[i] = pb.NewKVStoreClient(conn.clientConn[i])
    }
    return &KVStore{Conn: conn, client: clients}
}


func (c *KVStore) Put(k interface{}, v interface{}) error {

    ctx, cancel := context.WithTimeout(context.Background(), Deadline)
    defer cancel()
    replicaClients := c.selectNodes(k)
    if len(replicaClients) == 0 {
    	return fmt.Errorf("%v.Put(_) = _, %v: ", c.client, " no available nodes!")
	}
    // Iterate over replica client list and perform Put operation
    for _, client := range replicaClients  {
    	_, err := client.Put(ctx, &pb.KVPair{ Key: ToBytes(k), Value: ToBytes(v) }) 
    	if err != nil {
    		return fmt.Errorf("%v.Put(_) = _, %v: ", c.client, err)
		}
	}
	return nil
}


func (c *KVStore) splitBatch(batch map[interface{}]interface{}, replicas int) []map[interface{}]interface{} {

	c.Conn.nodeMapLock.RLock()
	defer c.Conn.nodeMapLock.RUnlock()

    batches := make([]map[interface{}]interface{}, len(c.client))
    for i, _ := range batches {
        batches[i] = make(map[interface{}]interface{}, 0)
    }
    for k, v := range batch {
        nodeKeys := c.Conn.hashTable.GetN(replicas, ToString(k))
/*
    	if len(nodeKeys) == 0 {
    		return fmt.Errorf("%v.splitBatch(_) = _, %v: key: %v", c.client, " no available nodes!", k)
		}
*/
    	// Iterate over node key list and collate into batches
    	for _, nodeKey := range nodeKeys  {
			if i, ok := c.Conn.nodeMap[nodeKey]; ok {
				batches[i][k] = v
			}
		}
	}
	return batches
}


func (c *KVStore) BatchPut(batch map[interface{}]interface{}) error {
  
    batches := c.splitBatch(batch, c.Replicas)

    done := make(chan error)
    defer close(done)
    count := len(batches)
    for i, v := range batches {
        go func(client pb.KVStoreClient, m map[interface{}]interface{}) {
			done <- c.batchPut(client, m)
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


func (c *KVStore) batchPut(client pb.KVStoreClient, batch map[interface{}]interface{}) error {

    ctx, cancel := context.WithTimeout(context.Background(), Deadline * time.Second)
    defer cancel()
    b := make([]*pb.KVPair, len(batch))
    i := 0
    for k, v := range batch {
	    b[i] = &pb.KVPair{ Key: ToBytes(k), Value: ToBytes(v) }
		i++
    }
    stream, err := client.BatchPut(ctx)
    if err != nil {
        return fmt.Errorf("%v.BatchPut(_) = _, %v: ", c.client, err)
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


func (c *KVStore) Lookup(key interface{}, valueType reflect.Kind) (interface{}, error) {

    ctx, cancel := context.WithTimeout(context.Background(), Deadline * time.Second)
    defer cancel()
    replicaClients := c.selectNodes(key)
    if len(replicaClients) == 0 {
        return nil, fmt.Errorf("%v.Lookup(_) = _, %v: ", c.client, " no available nodes!")
    }

    // Use the highest weight client
    lookup, err := replicaClients[0].Lookup(ctx, &pb.KVPair{ Key: ToBytes(key), Value: nil} ) 
    if err != nil {
        return uint64(0), fmt.Errorf("%v.Lookup(_) = _, %v: ", c.client, err)
    }
	if lookup.Value != nil {
    	return UnmarshalValue(valueType, lookup.Value), nil
	}
    return nil, nil
}


func (c *KVStore) BatchLookup(batch map[interface{}]interface{}) (map[interface{}]interface{}, error) {

    // We dont want to iterate over replicas for lookups so count is 1
    batches := c.splitBatch(batch, 1)

    results := make(map[interface{}]interface{}, 0)
    rchan := make(chan map[interface{}]interface{})
	defer close(rchan)
    done := make(chan error)
    defer close(done)
    count := len(batches)
    for i, _ := range batches {
        go func(client pb.KVStoreClient, b map[interface{}]interface{}) {
			r, e := c.batchLookup(client, b)
            rchan <- r
            done <- e
		}(c.client[i], batches[i])

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


func (c *KVStore) batchLookup(client pb.KVStoreClient, batch map[interface{}]interface{}) (map[interface{}]interface{},
		 error) {

    ctx, cancel := context.WithTimeout(context.Background(), Deadline * time.Second)
    defer cancel()
    stream, err := client.BatchLookup(ctx)
    if err != nil {
        return nil, fmt.Errorf("%v.BatchLookup(_) = _, %v", c.client, err)
    }

    var keyType, valueType reflect.Kind
	// Just grab first key/value of incoming data to determine types
    for k, v := range batch {
		keyType = reflect.ValueOf(k).Kind()
		valueType = reflect.ValueOf(v).Kind()
		break
	}

    waitc := make(chan struct{})
    results := make(map[interface{}]interface{}, len(batch))
    go func() {
        for {
            kv, err := stream.Recv()
            if err == io.EOF {
                // read done.
                close(waitc)
                return
            }
            if err != nil {
                c.err <- fmt.Errorf("Failed to receive a KV pair : %v", err)
				return
            }
			k := UnmarshalValue(keyType, kv.Key)
			v := UnmarshalValue(valueType, kv.Value)
            results[k] = v
        }
    }()
    for k, v := range batch {
        newKV := &pb.KVPair{Key: ToBytes(k), Value: ToBytes(v)}
        if err := stream.Send(newKV); err != nil {
            return nil, fmt.Errorf("Failed to send a KV pair: %v", err)
        }
    }
    stream.CloseSend()
    <-waitc
    return results, nil
}


func (c *KVStore) Items(keyType, valueType reflect.Kind) (map[interface{}]interface{}, error) {

    results := make(map[interface{}]interface{}, 0)
    rchan := make(chan map[interface{}]interface{})
    defer close(rchan)
    done := make(chan error)
    defer close(done)
    count := len(c.client)

    for i, _ := range c.client {
        go func(client pb.KVStoreClient) {
            r, e := c.items(client, keyType, valueType)
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


func (c *KVStore) items(client pb.KVStoreClient, keyType, valueType reflect.Kind) (map[interface{}]interface{}, error) {

    batch := make(map[interface{}]interface{}, 0)

    ctx, cancel := context.WithTimeout(context.Background(), Deadline)
    defer cancel()

    stream, err := client.Items(ctx, &empty.Empty{})
    if err != nil {
        return nil, fmt.Errorf("%v.Items(_) = _, %v", client, err)
    }
    for {
        item, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("%v.Items(_) = _, %v", client, err)
        }
		batch[UnmarshalValue(keyType, item.Key)] = UnmarshalValue(valueType, item.Value)
    }

    return batch, nil
}


func (c *KVStore) selectNodes(key interface{}) []pb.KVStoreClient {

	c.Conn.nodeMapLock.RLock()
	defer c.Conn.nodeMapLock.RUnlock()

    nodeKeys := c.Conn.hashTable.GetN(c.Conn.Replicas, ToString(key))
    selected := make([]pb.KVStoreClient, len(nodeKeys))
	
    for i, v := range nodeKeys {
    	if j, ok := c.Conn.nodeMap[v]; ok {
			selected[i] = c.client[j]
		}
	}
    
    return selected
}
