package server

import (
	"fmt"
	"log"
	"time"
	"github.com/hashicorp/consul/api"
	"github.com/stvp/rendezvous"
)

const (
	checkInterval = 5 * time.Second
	pollWait      = time.Second
)

// Node is a single node in a distributed hash table, coordinated using
// services registered in Consul. Key membership is determined using rendezvous
// hashing to ensure even distribution of keys and minimal key membership
// changes when a Node fails or otherwise leaves the hash table.
//
// Errors encountered when making blocking GET requests to the Consul agent API
// are logged using the log package.
type Node struct {
	// Consul
	serviceName string
	serviceID   string
	consul      *api.Client

	// health check endpoint
	checkURL      string
    Service       *EndPoint

	// Hash table
	hashTable *rendezvous.Table
	waitIndex uint64

	// Shutdown channels
	Stop chan bool
    Err  chan error
    
}

// Join creates a new Node and adds it to the distributed hash table specified
// by the given name. The given id should be unique among all Nodes in the hash
// table.
func Join(name string, e *EndPoint) (node *Node, err error) {
	node = &Node{
		serviceName: name,
		serviceID:   e.hashKey,
		Stop:        make(chan bool),
    	Err:		 make(chan error),
        checkURL:    "Status",
        Service:	 e,
	}

	node.consul, err = api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("node: can't create Consul API client: %s", err)
	}

    go func() {
		defer close(node.Stop)
    	node.Err <- node.Service.Start()
    }()
    
	err = node.register()
	if err != nil {
		return nil, fmt.Errorf("node: can't register %s service: %s", node.serviceName, err)
	}

	err = node.update()
	if err != nil {
		return nil, fmt.Errorf("node: can't fetch %s services list: %s", node.serviceName, err)
	}

	go node.poll()

	return node, nil
}

func (n *Node) register() (err error) {

	err = n.consul.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name: n.serviceName,
		ID:   n.serviceID,
		Check: &api.AgentServiceCheck{
            GRPC: fmt.Sprintf("%v:%v/%v", n.Service.BindAddr, n.Service.Port, n.checkURL),
			Interval: checkInterval.String(),
		},
	})
	return err
}

func (n *Node) poll() {
	var err error

	for {
		select {
		case <-n.Stop:
			return
        case err = <- n.Err:
			if err != nil {
				log.Printf("[node %s %s] error: %s", n.serviceName, n.serviceID, err)
			}
			return
		case <-time.After(pollWait):
			err = n.update()
			if err != nil {
				log.Printf("[node %s %s] error: %s", n.serviceName, n.serviceID, err)
			}
		}
	}
}

// update blocks until the service list changes or until the Consul agent's
// timeout is reached (10 minutes by default).
func (n *Node) update() (err error) {
	opts := &api.QueryOptions{WaitIndex: n.waitIndex}
	serviceEntries, meta, err := n.consul.Health().Service(n.serviceName, "", true, opts)
	if err != nil {
		return err
	}
    if serviceEntries == nil {
		return nil
	}

	ids := make([]string, len(serviceEntries))
	for i, entry := range serviceEntries {
		ids[i] = entry.Service.ID
	}

	n.hashTable = rendezvous.New(ids)
	n.waitIndex = meta.LastIndex

	return nil
}

// Member returns true if the given key belongs to this Node in the distributed
// hash table.
func (n *Node) Member(key string) bool {
	return n.hashTable.Get(key) == n.serviceID
}

// Leave removes the Node from the distributed hash table by de-registering it
// from Consul. Once Leave is called, the Node should be discarded. An error is
// returned if the Node is unable to successfully deregister itself from
// Consul. In that case, Consul's health check for the Node will fail and
// require manual cleanup.
func (n *Node) Leave() (err error) {
	close(n.Stop) // stop polling for state
	err = n.consul.Agent().ServiceDeregister(n.serviceID)
    //close(n.stop) FIXME: stop the server
	return err
}
