package dstash

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
    "github.com/RoaringBitmap/roaring"
    pb "github.com/guymolinari/dstash/grpc"
)

type BitmapIndex struct {
	*Conn
    client				[]pb.BitmapIndexClient
    batch		 	 	map[string]map[string]map[uint64]map[time.Time]*roaring.Bitmap
    batchSize         	int
    batchCount        	int
    batchMutex         	sync.Mutex
}


func NewBitmapIndex(conn *Conn, batchSize int) *BitmapIndex {

	clients := make([]pb.BitmapIndexClient, len(conn.clientConn))
	for i := 0; i < len(conn.clientConn); i++ {
		clients[i] = pb.NewBitmapIndexClient(conn.clientConn[i])
	}
	return &BitmapIndex{Conn: conn, batchSize: batchSize, client: clients}
}


func (c *BitmapIndex) Dispose() error {

	c.batchMutex.Lock()
    defer c.batchMutex.Unlock()

    if c.batch != nil {
        if err := c.BatchSetBit(c.batch); err != nil {
            return err
        }   
        c.batch = nil
    } 
	return nil
}


func (c *BitmapIndex) SetBit(index, field string, columnID, rowID uint64, ts time.Time) error {

    c.batchMutex.Lock()
    defer c.batchMutex.Unlock()

	if c.batch == nil {
		c.batch = make(map[string]map[string]map[uint64]map[time.Time]*roaring.Bitmap)
	}
	if _, ok := c.batch[index]; !ok {
		c.batch[index] = make(map[string]map[uint64]map[time.Time]*roaring.Bitmap)
	}
	if _, ok := c.batch[index][field]; !ok {
		c.batch[index][field] = make(map[uint64]map[time.Time]*roaring.Bitmap)
	}
	if _, ok := c.batch[index][field][rowID]; !ok {
		c.batch[index][field][rowID] = make(map[time.Time]*roaring.Bitmap)
	}
    if bmap, ok := c.batch[index][field][rowID][ts]; !ok {
		b := roaring.BitmapOf(uint32(columnID))
		c.batch[index][field][rowID][ts] = b
	} else {
		bmap.Add(uint32(columnID))
	}

    c.batchCount++

    if c.batchCount >= c.batchSize {

		//start := time.Now()
        if err := c.BatchSetBit(c.batch); err != nil {
            return err
        }
    	//elapsed := time.Since(start)
		//log.Printf("BatchBitSet duration %v\n", elapsed)
		//log.Printf("Batch create duration %v.\n", (c.batchDuration / float64(c.batchCount) * float64(1000000)))
        c.batch = nil
		c.batchCount = 0
    }
	return nil
}


func (c *BitmapIndex) BatchSetBit(batch map[string]map[string]map[uint64]map[time.Time]*roaring.Bitmap) error {

    batches := c.splitBitmapBatch(batch, c.Conn.Replicas)

    done := make(chan error)
    defer close(done)
    count := len(batches)
    for i, v := range batches {
        go func(client pb.BitmapIndexClient, m map[string]map[string]map[uint64]map[time.Time]*roaring.Bitmap) {
			done <- c.batchSetBit(client, m)
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


func (c *BitmapIndex) batchSetBit(client pb.BitmapIndexClient, 
		batch map[string]map[string]map[uint64]map[time.Time]*roaring.Bitmap) error {

    ctx, cancel := context.WithTimeout(context.Background(), Deadline)
    defer cancel()
    b := make([]*pb.IndexKVPair, 0)
    i := 0
   	for indexName, index := range batch {
        for fieldName, field := range index {
	        for rowID, ts := range field {
	        	for t, bitmap := range ts {
					buf, err := bitmap.ToBytes()
					if err != nil {
	        			log.Printf("bitmap.ToBytes: %v", err)
						return err
					}
					b = append(b, &pb.IndexKVPair{IndexPath: indexName + "/" + fieldName, 
						Key: ToBytes(rowID), Value: buf, Time: t.UnixNano()})
					i++
					//log.Printf("Sent batch %d for path %s\n", i, b[i].IndexPath)
				}
			}
		}
    }
    stream, err := client.BatchSetBit(ctx)
    if err != nil {
        log.Printf("%v.BatchSetBit(_) = _, %v: ", c.client, err)
        return fmt.Errorf("%v.BatchSetBit(_) = _, %v: ", c.client, err)
    }

	for i := 0; i < len(b); i++ {
        if err := stream.Send(b[i]); err != nil {
            log.Printf("%v.Send(%v) = %v", stream, b[i], err)
            return fmt.Errorf("%v.Send(%v) = %v", stream, b[i], err)
        }
	}
    _, err2 := stream.CloseAndRecv()
    if err2 != nil {
        log.Printf("%v.CloseAndRecv() got error %v, want %v", stream, err2, nil)
        return fmt.Errorf("%v.CloseAndRecv() got error %v, want %v", stream, err2, nil)
    }
    return nil
}


func (c *BitmapIndex) splitBitmapBatch(batch map[string]map[string]map[uint64]map[time.Time]*roaring.Bitmap, 
		replicas int) []map[string]map[string]map[uint64]map[time.Time]*roaring.Bitmap {

	c.Conn.nodeMapLock.RLock()
	defer c.Conn.nodeMapLock.RUnlock()

    batches := make([]map[string]map[string]map[uint64]map[time.Time]*roaring.Bitmap, len(c.client))
    for i, _ := range batches {
        batches[i] = make(map[string]map[string]map[uint64]map[time.Time]*roaring.Bitmap)
    }

	for indexName, index := range batch {
		for fieldName, field := range index {
			for rowID, ts := range field {
				nodeKeys := c.Conn.hashTable.GetN(replicas, fmt.Sprintf("%s/%s/%d", indexName, fieldName, rowID))
				for _, nodeKey := range nodeKeys {
					if i, ok := c.Conn.nodeMap[nodeKey]; ok {
						for t, bitmap := range ts {
							if batches[i] == nil {
								batches[i] = make(map[string]map[string]map[uint64]map[time.Time]*roaring.Bitmap)
							}
							if _, ok := batches[i][indexName]; !ok {
								batches[i][indexName] = make(map[string]map[uint64]map[time.Time]*roaring.Bitmap)
							}
							if _, ok := batches[i][indexName][fieldName]; !ok {
								batches[i][indexName][fieldName] = make(map[uint64]map[time.Time]*roaring.Bitmap)
							}
							if _, ok := batches[i][indexName][fieldName][rowID]; !ok {
								batches[i][indexName][fieldName][rowID] = make(map[time.Time]*roaring.Bitmap)
							}
							batches[i][indexName][fieldName][rowID][t] = bitmap
						}
					}
				}
			}
		}
	}
	return batches
}


func (c *BitmapIndex) query(query *pb.BitmapQuery) (*roaring.Bitmap, error) {

    ctx, cancel := context.WithTimeout(context.Background(), Deadline)
    defer cancel()

	c.Conn.nodeMapLock.RLock()
	defer c.Conn.nodeMapLock.RUnlock()

	// Split query fragments by shard key into separate client queries
    clientQueries := make([]*pb.BitmapQuery, len(c.client))
    for i, _ := range clientQueries {
       clientQueries[i] = &pb.BitmapQuery{Query: []*pb.QueryFragment{}, FromTime: query.FromTime, ToTime: query.ToTime}
    }

   	for _, v := range query.Query {
        if v.Index == "" {
            return nil, fmt.Errorf("Index not specified for query fragment %#v", v)
        }
        if v.Field == "" {
            return nil, fmt.Errorf("Field not specified for query fragment %#v", v)
        }
		nodeKeys := c.Conn.hashTable.GetN(c.Conn.Replicas, fmt.Sprintf("%s/%s/%d", v.Index, v.Field, v.RowID))
		for _, nodeKey := range nodeKeys {
    		if i, ok := c.nodeMap[nodeKey]; ok {
				clientQueries[i].Query = append(clientQueries[i].Query, v)
			}
		}
	}

	result := roaring.NewBitmap()
	for i, q := range clientQueries {
		if len(q.Query) == 0 {
			continue
		}
		if buf, err := c.client[i].Query(ctx, q); err != nil {
            return nil, fmt.Errorf("Error executing query - %v", err)
		} else {
			bm := roaring.NewBitmap()
			if err := bm.UnmarshalBinary(buf.Value); err != nil {
            	return nil, fmt.Errorf("Error unmarshalling query result - %v", err)
			} else {
				result = roaring.FastOr(result, bm)
			}
		}
	}
    return result, nil
}


func (c *BitmapIndex) CountQuery(query *pb.BitmapQuery) (uint64, error) {

	if result, err := c.query(query); err != nil {
		return 0, err
	} else {
		return result.GetCardinality(), nil
	}
}


func (c *BitmapIndex) ResultsQuery(query *pb.BitmapQuery, limit uint64) ([]uint32, error) {

	if result, err := c.query(query); err != nil {
		return []uint32{}, err
	} else {
		if result.GetCardinality() > limit {
			ra := make([]uint32, 0)
			iter := result.Iterator()
			for iter.HasNext() {
				ra = append(ra, iter.Next())
			}
			return ra, nil
		} else {
			return result.ToArray(), nil
		}
	}
}


