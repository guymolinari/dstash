package server

import (
    "bytes"
    "context"
    "encoding/binary"
	"hash"
	"hash/fnv"
    "fmt"
    "io"
	"io/ioutil"
    "log"
    "net"
    "os"
    "path"
	"path/filepath"
    "regexp"
	"sort"
	"strings"
	"strconv"
    "sync"
	"time"
    "golang.org/x/text/unicode/norm"
    "github.com/steakknife/bloomfilter"
    "github.com/RoaringBitmap/roaring"
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
	SEP = string(os.PathSeparator)
	
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
	bitmapCache	  map[string]map[string]map[uint64]*StandardBitmap
	bitmapCacheLock	  sync.RWMutex
	fragQueue	  chan *BitmapFragment
	workers		  int
	fragFileLock  sync.Mutex
	setBitThreads *CountTrigger
	writeSignal   chan bool
	
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

	m.dataDir = dataDir
	m.fragQueue = make(chan *BitmapFragment, 1000000)

	m.bitmapCache = make(map[string]map[string]map[uint64]*StandardBitmap)
	m.workers = 3
	m.writeSignal = make(chan bool, 1)
	m.setBitThreads = NewCountTrigger(m.writeSignal)

    return m, nil
}


func (m *EndPoint) Start() error {

    lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", m.BindAddr, m.Port))
    if err != nil {
        return fmt.Errorf("failed to listen: %v", err)
    }
    var opts []grpc.ServerOption
	opts = append(opts, grpc.MaxRecvMsgSize(1024*1024*20), grpc.MaxSendMsgSize(1024*1024*20))

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

	for i := 0; i < m.workers; i++ {
		go m.batchLoadProcessLoop(i + 1)
	}
	go m.overflowProcessLoop()

    // Read files from disk
	m.readBitmapFiles(m.fragQueue, true)

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

	m.setBitThreads.Add(1)
	defer m.setBitThreads.Add(-1)

    for {
        kv, err := stream.Recv()
        if err == io.EOF {
            return stream.SendAndClose(&empty.Empty{})
        }
        if err != nil {
            return err
        }
	    if kv == nil {
			return fmt.Errorf("KV Pair must not be nil")
		}
	    if kv.Key == nil || len(kv.Key) == 0 {
			return fmt.Errorf("Key must be specified.")
		}

        s := strings.Split(kv.IndexPath, "/")
		if len(s) != 2 {
			err = fmt.Errorf("IndexPath %s not valid.", kv.IndexPath)
			log.Printf("%s", err)
			return err
		}
		indexName := s[0]
		fieldName := s[1]
   	    rowID := binary.LittleEndian.Uint64(kv.Key)

		select {
		case m.fragQueue <- newBitmapFragment(indexName, fieldName, rowID, kv.Value):
		default:
			// Fragment queue is full, send to disk
			if err := m.journalBitmapFragment(kv.Value, indexName, fieldName, rowID); err != nil {
				log.Printf("%s", err)
				return err
			}
		}
    }
}


func (m *EndPoint) journalBitmapFragment(newBm []byte, indexName, fieldName string, rowID uint64) error {

	m.fragFileLock.Lock()
	defer m.fragFileLock.Unlock()

	if fd, err := m.newFragmentFile(indexName, fieldName, rowID); err == nil {
		if _, err := fd.Write(newBm); err != nil {
			return err
		}
		if err := fd.Close(); err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}


func (m *EndPoint) saveCompleteBitmap(bm *StandardBitmap, indexName, fieldName string, rowID uint64, create bool) error {

	data, err := bm.Bits.MarshalBinary() 
	if err != nil {
		return err
	}

	if fd, err := m.writeCompleteFile(indexName, fieldName, rowID, create); err == nil {
		if _, err := fd.Write(data); err != nil {
			return err
		}
		if err := fd.Close(); err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}


type StandardBitmap struct {
	Bits		*roaring.Bitmap
	ModTime		time.Time
	Lock		sync.RWMutex
}


func newStandardBitmap() *StandardBitmap {
	return &StandardBitmap{Bits: roaring.NewBitmap(), ModTime: time.Now().Add(time.Second * 10)}
}


type BitmapFragment struct {
	IndexName			string
	FieldName			string
	RowID				uint64
	BitData				[]byte
	ModTime				time.Time
}


func newBitmapFragment(index, field string, rowID uint64, f []byte) *BitmapFragment {
	return &BitmapFragment{IndexName: index, FieldName: field, RowID: rowID, BitData: f}
}


func (m *EndPoint) shouldWriteFile(index, field string, rowID uint64, mod time.Time) (create bool, update bool){

	baseDir := m.dataDir + SEP + "bitmap" + SEP + index + SEP + field
    os.MkdirAll(baseDir, 0755)
	fname := fmt.Sprintf("%d", rowID)
    path := baseDir + SEP + fname
   	info, err := os.Stat(path)
	if err != nil {
		create = true
		return
	}
	if mod.After(info.ModTime()) {
		update = true
		return
	}
	return
}


func (m *EndPoint) writeCompleteFile(index, field string, rowID uint64, create bool) (*os.File, error) {

	baseDir := m.dataDir + SEP + "bitmap" + SEP + index + SEP + field
    os.MkdirAll(baseDir, 0755)
	fname := fmt.Sprintf("%d", rowID)
    path := baseDir + SEP + fname
	if create {
		f, err := os.OpenFile(path, os.O_CREATE | os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		return f, nil
	} else {
		f, err := os.OpenFile(path, os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
}


func (m *EndPoint) newFragmentFile(index, field string, rowID uint64) (*os.File, error) {

	baseDir := m.dataDir + SEP + "overflow" + SEP + index + SEP + field + SEP +
		fmt.Sprintf("%d", rowID)
    os.MkdirAll(baseDir, 0755)

	replacer := strings.NewReplacer("-", "", ":", "", ".", "_")
	fname := replacer.Replace(time.Now().Format(time.RFC3339Nano))
    path := baseDir + SEP + fname
	f, err := os.OpenFile(path, os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return f, nil
}


func (m *EndPoint) readBitmapFiles(fragQueue chan *BitmapFragment, isComplete bool) error {

	m.fragFileLock.Lock()
	defer m.fragFileLock.Unlock()

	var subPath string
    if isComplete {
		subPath = "bitmap"
	} else {
		subPath = "overflow"
	}

	baseDir := m.dataDir + SEP + subPath
    os.MkdirAll(baseDir, 0755)

	var list []*BitmapFragment = make([]*BitmapFragment, 0)

	err := filepath.Walk(baseDir,
    	func(path string, info os.FileInfo, err error) error {
    	if err != nil {
        	return err
    	}
		if info.IsDir() {
			return nil
		}
		bf := &BitmapFragment{ModTime: info.ModTime()}

	    data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("readBitmapFiles: ioutil.ReadFile - %v", err)
			return err
		}
		s := strings.Split(path, SEP)
		if len(s) < 4 {
			err := fmt.Errorf("readBitmapFiles: Could not parse path [%s]", path)
			log.Println(err)
			return err
		}

		if isComplete {
			bf.IndexName = s[len(s) - 3]
			bf.FieldName = s[len(s) - 2]
			bf.RowID, err = strconv.ParseUint(s[len(s) - 1], 10, 64)
			if err != nil {
				err := fmt.Errorf("readBitmapFiles: Could not parse RowID - %v", err)
				log.Println(err)
				return err
			}
		} else {
			bf.IndexName = s[len(s) - 4]
			bf.FieldName = s[len(s) - 3]
			bf.RowID, err = strconv.ParseUint(s[len(s) - 2], 10, 64)
			if err != nil {
				err := fmt.Errorf("oldestFragmentFile: Could not parse RowID - %v", err)
				log.Println(err)
				return err
			}
	    	err = os.Remove(path)
			if err != nil {
				err := fmt.Errorf("readBitmapFiles: Could not delete file- %v", err)
				return err
			}
		}

		bf.BitData = data
		list = append(list, bf)

    	return nil
	})
	if err != nil {
    	log.Printf("filepath.Walk - %v", err)
		return err
	}

    sort.Slice(list, func(i, j int) bool { return list[i].ModTime.UnixNano() < list[j].ModTime.UnixNano() })
    if len(list) == 0 {
        return nil
    }

	for _, f := range list {
		fragQueue <- f
	}

	return nil
}


func (m *EndPoint) batchLoadProcessLoop(threadID int) {

	log.Printf("batchLoadProcess [Thread #%d] - Started.", threadID)
	for {
		select {
		case frag := <- m.fragQueue:
			m.updateBitmapCache(frag)
			continue
		default:
		}
		select {
		case frag := <- m.fragQueue:
			m.updateBitmapCache(frag)
			continue
		case <- m.writeSignal:
			m.checkPersistBitmapCache()
		}
	}
	log.Printf("batchLoadProcess [Thread #%d] - Ended.", threadID)
}


func (m *EndPoint) overflowProcessLoop() {

	log.Printf("overflowProcessLoop - Started.")
	for {
		if err := m.readBitmapFiles(m.fragQueue, false); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second)
	}
	log.Printf("overflowProcessLoop - Ended.")
}


func (m *EndPoint) updateBitmapCache(f *BitmapFragment) {

	// If the rowID exists then merge in the new set of bits

   	start := time.Now()
	newBm := newStandardBitmap()
	if err := newBm.Bits.UnmarshalBinary(f.BitData); err != nil {
		log.Printf("updateBitmapCache - UnmarshalBinary error - %v", err)
		return
	}
	m.bitmapCacheLock.Lock()
	if _, ok := m.bitmapCache[f.IndexName]; !ok {
		m.bitmapCache[f.IndexName] = make(map[string]map[uint64]*StandardBitmap)
	}
	if _, ok := m.bitmapCache[f.IndexName][f.FieldName]; !ok {
		m.bitmapCache[f.IndexName][f.FieldName] = make(map[uint64]*StandardBitmap)
	}
	if existBm, ok := m.bitmapCache[f.IndexName][f.FieldName][f.RowID]; !ok {
		m.bitmapCache[f.IndexName][f.FieldName][f.RowID] = newBm
		m.bitmapCacheLock.Unlock()
	} else {
		m.bitmapCacheLock.Unlock()
		existBm.Lock.Lock()
        existBm.Bits = roaring.ParOr(0, existBm.Bits, newBm.Bits)
		existBm.ModTime = time.Now()
		existBm.Lock.Unlock()
	}
    elapsed := time.Since(start)
	if elapsed.Nanoseconds() > 10000000 {
    	log.Printf("updateBitmapCache [%s/%s/%d] done in %v.\n", f.IndexName, f.FieldName, f.RowID,  elapsed)
	}
}


func (m *EndPoint) checkPersistBitmapCache() {

	m.bitmapCacheLock.RLock()
	defer m.bitmapCacheLock.RUnlock()

    for indexName, index := range m.bitmapCache {
        for fieldName, field := range index {
            for rowID, bitmap := range field {
				lastModSecs := time.Since(bitmap.ModTime).Round(time.Second)
				if dur, _ := time.ParseDuration("60s"); lastModSecs >= dur {
					bitmap.Lock.Lock()
   					start := time.Now()
					if create, update := m.shouldWriteFile(indexName, fieldName, rowID, bitmap.ModTime); create || update {
						if err := m.saveCompleteBitmap(bitmap, indexName, fieldName, rowID, create); err != nil {
							log.Printf("saveCompleteBitmap failed! - %v", err)
							bitmap.Lock.Unlock()
							return
						}
	    				elapsed := time.Since(start)
	    				log.Printf("Persist [%s/%s/%d] done in %v. Last mod %v seconds ago.", indexName, 
						 	fieldName, rowID,  elapsed, lastModSecs)
					}
					bitmap.Lock.Unlock()
				}
			}
		}
	}
}


func (m *EndPoint) Query(ctx context.Context, query *pb.BitmapQuery) (*wrappers.BytesValue, error) {

    if query == nil {
		return &wrappers.BytesValue{}, fmt.Errorf("query must not be nil")
	}
	
    if query.Query == nil {
		return &wrappers.BytesValue{}, fmt.Errorf("query fragment array must not be nil")
	}
    if len(query.Query) == 0 {
		return &wrappers.BytesValue{}, fmt.Errorf("query fragment array must not be empty")
	}
	prevOp := pb.QueryFragment_INTERSECT
	firstTime := true
	result := roaring.NewBitmap()
	for _, v := range query.Query {
		if v.Index == "" {
			return nil, fmt.Errorf("Index not specified for query fragment %#v", v)
		}
		if v.Field == "" {
			return nil, fmt.Errorf("Field not specified for query fragment %#v", v)
		}
		
		var bm *StandardBitmap 
		var ok bool
		m.bitmapCacheLock.RLock()
		defer m.bitmapCacheLock.RUnlock()
		if bm, ok = m.bitmapCache[v.Index][v.Field][v.RowID]; !ok {
			return &wrappers.BytesValue{Value: []byte("")}, 
				fmt.Errorf("Cannot find value for [%s/%s/%d]", v.Index, v.Field, v.RowID)
		}
		if firstTime  {
			prevOp = v.Operation
			result = bm.Bits
			firstTime = false
			continue
		}
		if prevOp == pb.QueryFragment_INTERSECT {
			result = roaring.ParAnd(0, result, bm.Bits)
		} else {
			result = roaring.ParOr(0, result, bm.Bits)
		}
		prevOp = v.Operation
	}

	if buf, err := result.MarshalBinary(); err != nil {
		return &wrappers.BytesValue{Value: []byte("")}, 
			fmt.Errorf("Cannot marshal result roaring bitmap - %v", err)
	} else {
		return &wrappers.BytesValue{Value: buf}, nil
	}
}


// Send message when counter reaches zero
type CountTrigger struct {
	num  int
	lock sync.Mutex
	trigger chan bool
}

func NewCountTrigger(t chan bool) *CountTrigger {
	return &CountTrigger{trigger: t}
}

/*
 * Add function provides thread safe addition of counter value based on input parameter.
 * If counter falls to zero then a value will be sent to trigger channel.
 */
func (c *CountTrigger) Add(n int) {
	c.lock.Lock()
	c.num += n
	if c.num == 0 {
		c.trigger <- true
	}
	c.lock.Unlock()
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

