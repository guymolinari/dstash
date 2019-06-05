# D-STASH - A simple high performance, resilient, distributed key/value store.
  
## Overview

DStash is a [GoLang](https://golang.org/) library and binary that implements a distributed hashtable that is persistent.  It utilizes a masterless architecture to achieve high availability and fault tolerance in support of applications that require high uptime SLAs. 

DStash builds upon the excellent open source library [Pogreb](https://github.com/akrylysov/pogreb).  Pogreb was designed to be deployed "in-process" with its consuming application.  Its interface is similar to Go's *map* data structure.  Persistence is implemented via *mmap* and its data set can be larger than available memory.  It is also designed to use relatively small memory footprint.

In addition to Pogreb, DStash leverages the following OSS projects:

- [Consul](https://github.com/hashicorp/consul) - For service discovery/health checking/locking and cluster management.
- [Grpc](https://grpc.io/docs) - To enable networking and service distribution.
- [Rendezvous](https://github.com/stvp/rendezvous) - A consistent hashing algorithm for mapping keys to server nodes.


## Use Cases

- Caching
- Simple key/value storage
- Data mapping
- Key management


## Why not just use Redis?

Redis is a very successful project and can fullfill all of the aforementioned use cases.  It does suffer from a few 
architectual flaws:

1. All data must fit in memory
2. Its HA/DR strategy is active/passive
3. Its clustering strategy is an afterthought.


It does offer a few advantages that DStash does not:

1. It supports other data structures such as queues and lists.
2. It supports TTL based expiration (although this is on the DStash roadmap).


## Performance Characteristics

Adding network serialization on top of Pogreb is sure to introduce latency.  Pogreb manages this by utilizing the streaming API constructs of Grpc.  API interfaces such as **Put()** and **Lookup()** also have batch equivalents that mitigates this when taken in aggregate.   So, a batch lookup of say 100,000 items can be performed in a few hundred milliseconds on contemporary hardware.  The same holds true for batch put operations.  

When DStash is deployed in the cloud, low cost server nodes allow for high levels of scalability.  The DStash client is optimized to use parallelism where ever possible.  DStash is architected to be deployed on an orchestrated platform such as Kubernetes to maximize administrative/operational scale.

//TODO:  Add performance metrics


## How to contribute

Pull requests are always welcome.  

//TODO: Add contributor's guide


## Roadmap

- Unit and integration tests sorely needed (this is a great way to get started as a contributor!)
- Examples/sample code/applications.
- Needs Makefile/build script for binary.
- Remove glide dependency in favor of Go Modules.
- Implement "scale in" with zero downtime on a running cluster.  Failover replication is already implemented but we need a "planned" way to reduce cluster size.
- Conversely we need to also "scale out" cluster size without downtime.
- Detailed performance benchmarks.
- Sample Kubernetes deployment scripts/best practices documentation.
- Dockerfile and entrypoint scripts.
- Remove dependency on Consul (native Kubernetes services?)
- Multi AZ/data center support.
- TTL based key expiration.
- Support for custom types.


## Node Deployment

```bash
usage: ./dstash [<flags>] [<data-dir>] [<bind>] [<port>]

DStash server node.

Flags:
  --help                 Show context-sensitive help (also try --help-long and --help-man).
  --tls                  Connection uses TLS if true.
  --cert-file=CERT-FILE  TLS cert file path.
  --key-file=KEY-FILE    TLS key file path.

Args:
  [<data-dir>]  Root directory for data files
  [<bind>]      Bind address for this endpoint.
  [<port>]      Port for this endpoint.
```

## Client API usage example

### Opening a connection

To open a connection, use the `dstash.NewClient()` and `dstash.OpenConnection()` functions:

```go
package main

import (
	"log"

	"github.com/guymolinari/dstash"
)

func main() {
    client := dstash.NewDefaultClient()
    if err := client.Connect(); err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect()
}
```

### Writing data
```go

// A key can be any value that implements the "Hash" interface, a value can be any native type
err := client.Put("myIndex", "myKey", "myValue")
if err != nil {
	log.Fatal(err)
}
```

### Reading data
```go

// A key can be any value that implements the "Hash" interface, a value can be any native type
val, err := client.Lookup("myIndex", "myKey")
if err != nil {
	log.Fatal(err)
}
```

