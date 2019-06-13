package dstash

import (
    "context"
    "encoding/binary"
    "io"
    "fmt"
    "log"
	"reflect"
    "sync"
    "time"
    "github.com/hashicorp/consul/api"
    "github.com/stvp/rendezvous"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/testdata"
    "github.com/golang/protobuf/ptypes/empty"
    "github.com/golang/protobuf/ptypes/wrappers"
    pb "github.com/guymolinari/dstash/grpc"
	"github.com/RoaringBitmap/roaring"
)

var (
    //Deadline int = 60
    Deadline time.Duration = time.Duration(60)
)

type Client struct {
    ServiceName			 string
    Quorum				 int
	Replicas			 int
    ConsulAgentAddr      string
    ServicePort			 int
    ServerHostOverride   string
    tls                  bool
    certFile             string
    client               []pb.DStashClient
    clientConn           []*grpc.ClientConn
    err                  chan error
    stop				 chan struct{}
	consul      	 	 *api.Client
    hashTable 			 *rendezvous.Table
    waitIndex 			 uint64
    pollWait             time.Duration
    nodes                []*api.ServiceEntry
    nodeMap              map[string]int
	nodeMapLock			 sync.RWMutex
	indexBatch           map[string]struct{}
    batchSize            int
    indexMutex           sync.Mutex
    bitBatch		 	 map[string]map[string]map[uint64]*roaring.Bitmap
    bitBatchSize         int
    bitBatchCount        int
    bitBatchMutex        sync.Mutex
	bitBatchDuration	 float64
}


func NewDefaultClient() *Client {
    m := &Client{}
	m.ServiceName = "dstash"
    m.ServicePort = 5000
    m.ServerHostOverride = "x.test.youtube.com"
    m.pollWait = time.Second * 5
    m.Quorum = 1
	m.Replicas = 1
    m.batchSize = 1000
	m.bitBatchSize = 1000000
    return m
}


func (m *Client) Connect() (err error) {

    m.consul, err = api.NewClient(api.DefaultConfig())
    if err != nil {
        return fmt.Errorf("client: can't create Consul API client: %s", err)
    }

    err = m.update()
    if err != nil {
        return fmt.Errorf("node: can't fetch %s services list: %s", m.ServiceName, err)
    }

    if m.hashTable == nil {
        return fmt.Errorf("node: uninitialized %s services list: %s", m.ServiceName, err)
    }

    if len(m.nodes) < m.Quorum {
        return fmt.Errorf("node: quorum size %d currently,  must be at least %d to handle requests for service %s", 
			len(m.nodes), m.Quorum, m.ServiceName)
    }

    m.client = make([]pb.DStashClient, len(m.nodes))
    m.clientConn = make([]*grpc.ClientConn, len(m.nodes))
   
	for i := 0; i < len(m.nodes); i++ {
	    var opts []grpc.DialOption
	    if m.tls {
	        if m.certFile == "" {
	            m.certFile = testdata.Path("ca.pem")
	        }
	        creds, err2 := credentials.NewClientTLSFromFile(m.certFile, m.ServerHostOverride)
	        if err2 != nil {
				err = fmt.Errorf("Failed to create TLS credentials %v", err2)
				return
	        }
	        opts = append(opts, grpc.WithTransportCredentials(creds))
	    } else {
	        opts = append(opts, grpc.WithInsecure())
	    }
	    m.clientConn[i], err = grpc.Dial(fmt.Sprintf("%s:%d", m.nodes[i].Node.Address, m.ServicePort), opts...)
	    if err != nil {
	        log.Fatalf("fail to dial: %v", err)
	    }
	
	    m.client[i] = pb.NewDStashClient(m.clientConn[i])
	}
	
    m.stop = make(chan struct{})
    go m.poll()
    return
}


func (c *Client) selectNodes(key interface{}) []pb.DStashClient {

	c.nodeMapLock.RLock()
	defer c.nodeMapLock.RUnlock()

    nodeKeys := c.hashTable.GetN(c.Replicas, ToString(key))
    selected := make([]pb.DStashClient, len(nodeKeys))
	
    for i, v := range nodeKeys {
    	if j, ok := c.nodeMap[v]; ok {
			selected[i] = c.client[j]
		}
	}
    
    return selected
}


func (c *Client) Disconnect() error {

    c.indexMutex.Lock()
    defer c.indexMutex.Unlock()

    if c.indexBatch != nil {
		if err := c.BatchIndex(c.indexBatch); err != nil {
			return err
		}
		c.indexBatch = nil
    }

    if c.bitBatch != nil {
		if err := c.BatchSetBit(c.bitBatch); err != nil {
			return err
		}
		c.bitBatch = nil
    }

	for i := 0; i < len(c.clientConn); i++ {
	    if c.clientConn[i] == nil {
			return fmt.Errorf("client: error - attempt to close nil clientConn.")
		}
	    c.clientConn[i].Close()
    	c.client[i] = nil
	}
    close(c.stop)
	return nil
}


func printStatus(client pb.DStashClient) {
    fmt.Print("Checking status ... ")
    ctx, cancel := context.WithTimeout(context.Background(), Deadline * time.Second)
    defer cancel()
    status, err := client.Status(ctx, &empty.Empty{}) 
    if err != nil {
        log.Fatalf("%v.Status(_) = _, %v: ", client, err)
    }
    fmt.Println(status)
}


func (c *Client) Put(k interface{}, v interface{}) error {

    ctx, cancel := context.WithTimeout(context.Background(), Deadline * time.Second)
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


func (c *Client) splitBatch(batch map[interface{}]interface{}, replicas int) []map[interface{}]interface{} {

	c.nodeMapLock.RLock()
	defer c.nodeMapLock.RUnlock()

    batches := make([]map[interface{}]interface{}, len(c.client))
    for i, _ := range batches {
        batches[i] = make(map[interface{}]interface{}, 0)
    }
    for k, v := range batch {
        nodeKeys := c.hashTable.GetN(replicas, ToString(k))
/*
    	if len(nodeKeys) == 0 {
    		return fmt.Errorf("%v.splitBatch(_) = _, %v: key: %v", c.client, " no available nodes!", k)
		}
*/
    	// Iterate over node key list and collate into batches
    	for _, nodeKey := range nodeKeys  {
			if i, ok := c.nodeMap[nodeKey]; ok {
				batches[i][k] = v
			}
		}
	}
	return batches
}


func (c *Client) splitStringBatch(batch map[string]struct{}, replicas int) []map[string]struct{} {

	c.nodeMapLock.RLock()
	defer c.nodeMapLock.RUnlock()

    batches := make([]map[string]struct{}, len(c.client))
    for i, _ := range batches {
        batches[i] = make(map[string]struct{}, 0)
    }
    for k, v := range batch {
        nodeKeys := c.hashTable.GetN(replicas, ToString(k))
    	// Iterate over node key list and collate into batches
    	for _, nodeKey := range nodeKeys  {
			if i, ok := c.nodeMap[nodeKey]; ok {
				batches[i][k] = v
			}
		}
	}
	return batches
}



func (c *Client) BatchPut(batch map[interface{}]interface{}) error {
  
    batches := c.splitBatch(batch, c.Replicas)

    done := make(chan error)
    defer close(done)
    count := len(batches)
    for i, v := range batches {
        go func(client pb.DStashClient, m map[interface{}]interface{}) {
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


func (c *Client) batchPut(client pb.DStashClient, batch map[interface{}]interface{}) error {

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


func (c *Client) Lookup(key interface{}, valueType reflect.Kind) (interface{}, error) {

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


func (c *Client) BatchLookup(batch map[interface{}]interface{}) (map[interface{}]interface{}, error) {

    // We dont want to iterate over replicas for lookups so count is 1
    batches := c.splitBatch(batch, 1)

    results := make(map[interface{}]interface{}, 0)
    rchan := make(chan map[interface{}]interface{})
	defer close(rchan)
    done := make(chan error)
    defer close(done)
    count := len(batches)
    for i, _ := range batches {
        go func(client pb.DStashClient, b map[interface{}]interface{}) {
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


func (c *Client) batchLookup(client pb.DStashClient, batch map[interface{}]interface{}) (map[interface{}]interface{},
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


func (c *Client) Items(keyType, valueType reflect.Kind) (map[interface{}]interface{}, error) {

    results := make(map[interface{}]interface{}, 0)
    rchan := make(chan map[interface{}]interface{})
    defer close(rchan)
    done := make(chan error)
    defer close(done)
    count := len(c.client)

    for i, _ := range c.client {
        go func(client pb.DStashClient) {
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


func (c *Client) items(client pb.DStashClient, keyType, valueType reflect.Kind) (map[interface{}]interface{}, error) {

    batch := make(map[interface{}]interface{}, 0)

    ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Second)
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


func (c *Client) Index(str string) error {
  
    c.indexMutex.Lock()
    defer c.indexMutex.Unlock()

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
   

func (c *Client) BatchIndex(batch map[string]struct{}) error {

    batches := c.splitStringBatch(batch, c.Replicas)

    done := make(chan error)
    defer close(done)
    count := len(batches)
    for i, v := range batches {
        go func(client pb.DStashClient, m map[string]struct{}) {
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


func (c *Client) batchIndex(client pb.DStashClient, batch map[string]struct{}) error {

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


func (c *Client) Search(searchTerms string) (map[uint64]struct{}, error) {

    results := make(map[uint64]struct{}, 0)
    rchan := make(chan map[uint64]struct{})
    defer close(rchan)
    done := make(chan error)
    defer close(done)
    count := len(c.client)

    for i, _ := range c.client {
        go func(client pb.DStashClient) {
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


func (c *Client) search(client pb.DStashClient, searchTerms string) (map[uint64]struct{}, error) {

    batch := make(map[uint64]struct{}, 0)

    ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Second)
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


func (c *Client) SetBit(index, field string, columnID, rowID uint64, ts time.Time) error {

    c.bitBatchMutex.Lock()
    defer c.bitBatchMutex.Unlock()

	if c.bitBatch == nil {
		c.bitBatch = make(map[string]map[string]map[uint64]*roaring.Bitmap)
	}
	if _, ok := c.bitBatch[index]; !ok {
		c.bitBatch[index] = make(map[string]map[uint64]*roaring.Bitmap)
	}
	if _, ok := c.bitBatch[index][field]; !ok {
		c.bitBatch[index][field] = make(map[uint64]*roaring.Bitmap)
	}
    if bmap, ok := c.bitBatch[index][field][rowID]; !ok {
		b := roaring.BitmapOf(uint32(columnID))
		c.bitBatch[index][field][rowID] = b
	} else {
		bmap.Add(uint32(columnID))
	}

    c.bitBatchCount++

    //c.bitBatchDuration = c.bitBatchDuration + elapsed.Seconds()
    if c.bitBatchCount >= c.bitBatchSize {
		//start := time.Now()
        if err := c.BatchSetBit(c.bitBatch); err != nil {
            return err
        }
    	//elapsed := time.Since(start)
		//log.Printf("BatchBitSet duration %v\n", elapsed)
		//log.Printf("Batch create duration %v.\n", (c.bitBatchDuration / float64(c.bitBatchCount) * float64(1000000)))
        c.bitBatch = nil
		c.bitBatchCount = 0
		//c.bitBatchDuration = float64(0)
    }
	return nil
}


func (c *Client) BatchSetBit(batch map[string]map[string]map[uint64]*roaring.Bitmap) error {

    batches := c.splitBitmapBatch(batch, c.Replicas)

    done := make(chan error)
    defer close(done)
    count := len(batches)
    for i, v := range batches {
        go func(client pb.DStashClient, m map[string]map[string]map[uint64]*roaring.Bitmap) {
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


func (c *Client) batchSetBit(client pb.DStashClient, batch map[string]map[string]map[uint64]*roaring.Bitmap) error {

    ctx, cancel := context.WithTimeout(context.Background(), Deadline * time.Second)
    defer cancel()
    b := make([]*pb.IndexKVPair, 0)
    i := 0
   	for indexName, index := range batch {
        for fieldName, field := range index {
            for rowID, bitmap := range field {
				buf, err := bitmap.ToBytes()
				if err != nil {
        			log.Printf("bitmap.ToBytes: %v", err)
					return err
				}
				b = append(b, &pb.IndexKVPair{IndexPath: indexName + "/" + fieldName, Key: ToBytes(rowID), Value: buf})
				i++
				//log.Printf("Sent batch %d for path %s\n", i, b[i].IndexPath)
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


func (c *Client) splitBitmapBatch(batch map[string]map[string]map[uint64]*roaring.Bitmap, 
		replicas int) []map[string]map[string]map[uint64]*roaring.Bitmap {

	c.nodeMapLock.RLock()
	defer c.nodeMapLock.RUnlock()

    batches := make([]map[string]map[string]map[uint64]*roaring.Bitmap, len(c.client))
    for i, _ := range batches {
        batches[i] = make(map[string]map[string]map[uint64]*roaring.Bitmap)
    }

	for indexName, index := range batch {
		for fieldName, field := range index {
			for rowID, bitmap := range field {
				nodeKeys := c.hashTable.GetN(replicas, fmt.Sprintf("%s/%s/%d", indexName, fieldName, rowID))
				for _, nodeKey := range nodeKeys {
					if i, ok := c.nodeMap[nodeKey]; ok {
						if batches[i] == nil {
							batches[i] = make(map[string]map[string]map[uint64]*roaring.Bitmap)
						}
						if _, ok := batches[i][indexName]; !ok {
							batches[i][indexName] = make(map[string]map[uint64]*roaring.Bitmap)
						}
						if _, ok := batches[i][indexName][fieldName]; !ok {
							batches[i][indexName][fieldName] = make(map[uint64]*roaring.Bitmap)
						}
						batches[i][indexName][fieldName][rowID] = bitmap
					}
				}
			}
		}
	}
	return batches
}


func (c *Client) poll() {
	var err error

	for {
		select {
		case <-c.stop:
			return
        case err = <- c.err:
			if err != nil {
				log.Printf("[client %s] error: %s", c.ServiceName, err)
			}
			return
		case <-time.After(c.pollWait):
			err = c.update()
			if err != nil {
				log.Printf("[client %s] error: %s", c.ServiceName, err)
			}
		}
	}
}

// update blocks until the service list changes or until the Consul agent's
// timeout is reached (10 minutes by default).
func (c *Client) update() (err error) {


	opts := &api.QueryOptions{WaitIndex: c.waitIndex}
	serviceEntries, meta, err := c.consul.Health().Service(c.ServiceName, "", true, opts)
	if err != nil {
		return err
	}
    if serviceEntries == nil {
		return nil
	}

	c.nodeMapLock.Lock()
	defer c.nodeMapLock.Unlock()

    c.nodes = serviceEntries
    c.nodeMap = make(map[string]int, len(serviceEntries))

	ids := make([]string, len(serviceEntries))
	for i, entry := range serviceEntries {
		ids[i] = entry.Service.ID
	    c.nodeMap[entry.Service.ID] = i	
	}

	c.hashTable = rendezvous.New(ids)
	c.waitIndex = meta.LastIndex

	return nil
}


func ToString(v interface{}) string {

    switch v.(type) {
	case string:
		return v.(string)
	default:
		return fmt.Sprintf("%v", v)
	}
}


func ToBytes(v interface{}) []byte {

    switch v.(type) {
	case string:
		return []byte(v.(string))
	case uint64:
    	b := make([]byte, 8)
    	binary.LittleEndian.PutUint64(b, v.(uint64))
		return b
	}
	panic("Should not be here")
}


func UnmarshalValue(kind reflect.Kind, buf []byte) interface{} {

	switch kind {
	case reflect.String:
		return string(buf)
	case reflect.Uint64:
		return binary.LittleEndian.Uint64(buf)
		
	}
	msg := fmt.Sprintf("Should not be here for kind [%s]!", kind.String())
    panic(msg)
}

