package server

import (
    "bytes"
    "context"
    "encoding/binary"
	"hash"
	"hash/fnv"
    "fmt"
    "io"
    "log"
    "net"
    "path"
    "regexp"
    "sync"
    "golang.org/x/text/unicode/norm"
    "github.com/steakknife/bloomfilter"
    "github.com/aviddiviner/go-murmur"
    "github.com/akrylysov/pogreb"
    "github.com/bbalet/stopwords"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/testdata"
    "github.com/golang/protobuf/ptypes/empty"
    "github.com/golang/protobuf/ptypes/wrappers"
    "google.golang.org/grpc/health/grpc_health_v1"
    pb "github.com/guymolinari/dstash/grpc"
)

var (
	//wordSegmenter = regexp.MustCompile(`[\pL\p{Mc}\p{Mn}-_']+`)
    wordSegmenter = regexp.MustCompile(`[\pL\p{Mc}\p{Mn}\p{Nd}-_']+`)
)

const (
  maxElements = 100
  probCollide = 0.0000001
)

type EndPoint struct {
    Port 		  uint
    BindAddr      string
    tls           bool
    certFile      string
    keyFile       string
	dataDir		  string
    store         *pogreb.DB
    searchIndex   *pogreb.DB
    hashKey       string
	storeCache	  map[string]*pogreb.DB
	storeLock	  sync.RWMutex
}


func NewEndPoint(dataDir string) (*EndPoint, error) {
    db, err := pogreb.Open(dataDir + "/" + "default.dat", nil)
    if err != nil {
        return nil, err
    }
    searchdb, err := pogreb.Open(dataDir + "/" + "search.dat", nil)
    if err != nil {
        return nil, err
    }
    m := &EndPoint{store: db, searchIndex: searchdb}
    m.hashKey = path.Base(dataDir)  // leaf directory name is consisten hash key
    if m.hashKey == "" || m.hashKey == "/" {
        return nil, fmt.Errorf("Data dir must not be root.")
	}

	m.storeCache = make(map[string]*pogreb.DB)
	m.dataDir = dataDir

    return m, nil
}


func (m *EndPoint) getStoreHandle(dataDir, indexPath string) (*pogreb.DB, error) {

    path := dataDir + "/" + indexPath
	if db, ok := m.storeCache[path]; ok {
		return db, nil
	}
    db, err := pogreb.Open(path, nil)
    if err != nil {
        return nil, err
    }
	m.storeCache[path] = db
	return db, nil
}


func (m *EndPoint) syncStores()  {
	
	for _, db := range m.storeCache {
		db.Sync()
	}
}


func (m *EndPoint) Start() error {

    lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", m.BindAddr, m.Port))
    if err != nil {
        return fmt.Errorf("failed to listen: %v", err)
    }
    var opts []grpc.ServerOption
    if m.tls {
        if m.certFile == "" {
            m.certFile = testdata.Path("server1.pem")
        }
        if m.keyFile == "" {
            m.keyFile = testdata.Path("server1.key")
        }
        creds, err := credentials.NewServerTLSFromFile(m.certFile, m.keyFile)
        if err != nil {
            return fmt.Errorf("Failed to generate credentials %v", err)
        }
        opts = []grpc.ServerOption{grpc.Creds(creds)}
    }
    grpcServer := grpc.NewServer(opts...)
    pb.RegisterDStashServer(grpcServer, m)
    grpc_health_v1.RegisterHealthServer(grpcServer, &HealthImpl{})
    grpcServer.Serve(lis)
    return nil
}


// Put - Insert a new key
func (m *EndPoint) Put(ctx context.Context, kv *pb.KVPair) (*empty.Empty, error) {
    if kv == nil {
		return &empty.Empty{}, fmt.Errorf("KV Pair must not be nil")
	}
    if kv.Key == nil || len(kv.Key) == 0 {
		return &empty.Empty{}, fmt.Errorf("Key must be specified.")
	}
	err := m.store.Put(kv.Key, kv.Value)
	if err != nil {
    	return &empty.Empty{}, err
	}
    return &empty.Empty{}, nil
}


func (m *EndPoint) Status(ctx context.Context, e *empty.Empty) (*pb.StatusMessage, error) {
    log.Printf("Status ping returning OK.\n")
    return &pb.StatusMessage{Status: "OK"}, nil
}


func (m *EndPoint) Lookup(ctx context.Context, kv *pb.KVPair) (*pb.KVPair, error) {
    if kv == nil {
		return &pb.KVPair{}, fmt.Errorf("KV Pair must not be nil")
	}
    if kv.Key == nil || len(kv.Key) == 0 {
		return &pb.KVPair{}, fmt.Errorf("Key must be specified.")
	}
	val, err := m.store.Get(kv.Key)
	if err != nil {
    	b := make([]byte, 8)
    	binary.LittleEndian.PutUint64(b, 0)
        kv.Value = b
		return kv, err
	} else {
        kv.Value = val
	}
	return kv, nil
}


func (m *EndPoint) BatchPut(stream pb.DStash_BatchPutServer) error {

    var putCount int32
    for {
        kv, err := stream.Recv()
        if err == io.EOF {
            m.store.Sync()
            return stream.SendAndClose(&empty.Empty{})
        }
        if err != nil {
            return err
        }
        putCount++
	    if kv == nil {
			return fmt.Errorf("KV Pair must not be nil")
		}
	    if kv.Key == nil || len(kv.Key) == 0 {
			return fmt.Errorf("Key must be specified.")
		}
		if err := m.store.Put(kv.Key, kv.Value); err != nil {
    		return err
		}
    }
}


func (m *EndPoint) BatchLookup(stream pb.DStash_BatchLookupServer) error {

    for {
        kv, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
		val, err := m.store.Get(kv.Key)
		if err != nil {
    		b := make([]byte, 8)
    		binary.LittleEndian.PutUint64(b, 0)
   	     	kv.Value = b
			return err
		} else {
   	     	kv.Value = val
		}
        if err := stream.Send(kv); err != nil {
            return err
        }
    }
}


func (m *EndPoint) Items(e *empty.Empty, stream pb.DStash_ItemsServer) error {

	it := m.store.Items()
	for {
	    key, val, err := it.Next()
	    if err != nil {
	        if err != pogreb.ErrIterationDone {
	            return err
	        }
	        break
	    }
        if err := stream.Send(&pb.KVPair{Key: key, Value: val}); err != nil {
            return err
        }
	}
    return nil
}


func (m *EndPoint) BatchIndex(stream pb.DStash_BatchIndexServer) error {

    for {
        sv, err := stream.Recv()
        if err == io.EOF {
            m.searchIndex.Sync()
            return stream.SendAndClose(&empty.Empty{})
        }
        if err != nil {
            return err
        }
        str := sv.GetValue()
	    if sv == nil || str == "" {
			return fmt.Errorf("String value must not be empty")
		}

        // Key is hash of original string
 		hashVal := uint64(murmur.MurmurHash2([]byte(str), 0))
    	key := make([]byte, 8)
    	binary.LittleEndian.PutUint64(key, hashVal)

		if found, err := m.searchIndex.Has(key); err != nil {
    		return err
		} else if found {
			continue
		}

        // Construct bloom filter sans stopwords
		bloomFilter, err := constructBloomFilter(str)
		if err != nil {
  			return err
		}

        bfBuf, err := bloomFilter.MarshalBinary()
		if err != nil {
  			return err
		}
        
		if err := m.searchIndex.Put(key, bfBuf); err != nil {
    		return err
		}
    }
}


func (m *EndPoint) Search(searchStr *wrappers.StringValue, stream pb.DStash_SearchServer) error {

    search := searchStr.GetValue()
	if searchStr == nil || search == "" {
		return fmt.Errorf("Search string must not be empty")
	}
    terms := parseTerms(search)

    hashedTerms := make([]hash.Hash64, len(terms))
    
    for i, v := range terms {
        hasher := fnv.New64a()
        hasher.Write(v)
		hashedTerms[i] = hasher
    }

	bloomFilter, err := bloomfilter.NewOptimal(maxElements, probCollide)
	if err != nil {
  		return err
	}

	it := m.searchIndex.Items()
	Top: for {
	    stringHash, val, err := it.Next()
	    if err != nil {
	        if err != pogreb.ErrIterationDone {
	            return err
	        }
	        break
	    }

        bloomFilter.UnmarshalBinary(val)

        // Perform "and" comparison. Item will be selected if all terms are contained.
        for _, v := range hashedTerms {
            if ! bloomFilter.Contains(v) {
				continue Top
			}
		}

        // return the hash of the original string value
   	    v := binary.LittleEndian.Uint64(stringHash[:8])
        if err := stream.Send(&wrappers.UInt64Value{Value: v}); err != nil {
            return err
        }
	}
	return nil
}


func parseTerms(content string) [][]byte {

    cleanStr := stopwords.CleanString(content, "en", true)
    c := norm.NFC.Bytes([]byte(cleanStr))
	c = bytes.ToLower(c)
	return wordSegmenter.FindAll(c, -1)
}


func constructBloomFilter(content string) (*bloomfilter.Filter, error) {

	words :=  parseTerms(content)

    // Construct bloom filter sans stopwords
	//bloomFilter, err := bloomfilter.NewOptimal(uint64(len(words)), probCollide)
	bloomFilter, err := bloomfilter.NewOptimal(uint64(maxElements), probCollide)
	if err != nil {
  		return nil, err
	}

    for _, v := range words {
        hasher := fnv.New64a()
		hasher.Write(v)
		bloomFilter.Add(hasher)
	}
    return bloomFilter, nil
}
        

func (m *EndPoint) BatchSetBit(stream pb.DStash_BatchSetBitServer) error {

    var putCount int32
    for {
        kv, err := stream.Recv()
        if err == io.EOF {
            m.syncStores()
            return stream.SendAndClose(&empty.Empty{})
        }
        if err != nil {
            return err
        }
        putCount++
	    if kv == nil {
			return fmt.Errorf("KV Pair must not be nil")
		}
	    if kv.Key == nil || len(kv.Key) == 0 {
			return fmt.Errorf("Key must be specified.")
		}
		store, err := m.getStoreHandle(m.dataDir, kv.IndexPath)
        if err != nil {
            return err
        }
		if err := store.Put(kv.Key, kv.Value); err != nil {
    		return err
		}
    }
}

type HealthImpl struct{}

// Check implements the health check interface, which directly returns to health status. There are also more complex health check strategies, such as returning based on server load.
func (h *HealthImpl) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
    //log.Printf("Health checking ...\n")
    return &grpc_health_v1.HealthCheckResponse{
        Status: grpc_health_v1.HealthCheckResponse_SERVING,
    }, nil
}

func (h *HealthImpl) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
    return nil
}

