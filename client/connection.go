package dstash

import (
    "context"
    "encoding/binary"
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
    pb "github.com/guymolinari/dstash/grpc"
)

var (
    Deadline time.Duration = time.Duration(500) * time.Second
)

type Conn struct {
    ServiceName			 string
    Quorum				 int
	Replicas			 int
    ConsulAgentAddr      string
    ServicePort			 int
    ServerHostOverride   string
    tls                  bool
    certFile             string
    admin                []pb.ClusterAdminClient
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
}


func NewDefaultConnection() *Conn {
    m := &Conn{}
	m.ServiceName = "dstash"
    m.ServicePort = 5000
    m.ServerHostOverride = "x.test.youtube.com"
    m.pollWait = time.Second * 5
    m.Quorum = 1
	m.Replicas = 1
    return m
}


func (m *Conn) Connect() (err error) {

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

    m.clientConn = make([]*grpc.ClientConn, len(m.nodes))
    m.admin = make([]pb.ClusterAdminClient, len(m.nodes))
	maxMsgSize := 1024*1024*20
   
	for i := 0; i < len(m.nodes); i++ {
	    var opts []grpc.DialOption
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize), 
			grpc.MaxCallSendMsgSize(maxMsgSize)))
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
		m.admin[i] = pb.NewClusterAdminClient(m.clientConn[i])
	
	}
	
    m.stop = make(chan struct{})
    go m.poll()
    return
}


func (c *Conn) selectNodes(key interface{}) []pb.ClusterAdminClient {

	c.nodeMapLock.RLock()
	defer c.nodeMapLock.RUnlock()

    nodeKeys := c.hashTable.GetN(c.Replicas, ToString(key))
    selected := make([]pb.ClusterAdminClient, len(nodeKeys))
	
    for i, v := range nodeKeys {
    	if j, ok := c.nodeMap[v]; ok {
			selected[i] = c.admin[j]
		}
	}
    
    return selected
}


func (c *Conn) Disconnect() error {

	for i := 0; i < len(c.clientConn); i++ {
	    if c.clientConn[i] == nil {
			return fmt.Errorf("client: error - attempt to close nil clientConn.")
		}
	    c.clientConn[i].Close()
    	c.admin[i] = nil
	}
	select {
	case _, open := <- c.stop:
		if open {
    		close(c.stop)
		}
	default:
		close(c.stop)
	}
	return nil
}


func printStatus(client pb.ClusterAdminClient) {
    fmt.Print("Checking status ... ")
    ctx, cancel := context.WithTimeout(context.Background(), Deadline * time.Second)
    defer cancel()
    status, err := client.Status(ctx, &empty.Empty{}) 
    if err != nil {
        log.Fatalf("%v.Status(_) = _, %v: ", client, err)
    }
    fmt.Println(status)
}


func (c *Conn) poll() {
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
func (c *Conn) update() (err error) {


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

