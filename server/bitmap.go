package server

import (
    "context"
    "encoding/binary"
    "fmt"
    "io"
	"io/ioutil"
    "log"
    "os"
	"path/filepath"
	"sort"
	"strings"
	"strconv"
    "sync"
	"time"
    "github.com/RoaringBitmap/roaring"
    "github.com/golang/protobuf/ptypes/empty"
    "github.com/golang/protobuf/ptypes/wrappers"
    "github.com/guymolinari/dstash/shared"
    pb "github.com/guymolinari/dstash/grpc"
)

type BitmapIndex struct {
    *EndPoint
    bitmapCache   		map[string]map[string]map[uint64]map[time.Time]*StandardBitmap
    bitmapCacheLock   	sync.RWMutex
    fragQueue     		chan *BitmapFragment
    workers       		int
    fragFileLock  		sync.Mutex
    setBitThreads 		*CountTrigger
    writeSignal   		chan bool
	tableCache	  		map[string]*shared.Table
	tableCacheLock		sync.RWMutex
}


func NewBitmapIndex(endPoint *EndPoint) *BitmapIndex {

    e := &BitmapIndex{EndPoint: endPoint}
	e.tableCache = make(map[string]*shared.Table)
    pb.RegisterBitmapIndexServer(endPoint.server, e)
    return e
}


func (m *BitmapIndex) Init() error {

    m.fragQueue = make(chan *BitmapFragment, 1000000)
    m.bitmapCache = make(map[string]map[string]map[uint64]map[time.Time]*StandardBitmap)
    m.workers = 3
    m.writeSignal = make(chan bool, 1)
    m.setBitThreads = NewCountTrigger(m.writeSignal)

    for i := 0; i < m.workers; i++ {
        go m.batchLoadProcessLoop(i + 1)
    }
    go m.overflowProcessLoop()

    // Read files from disk
    m.readBitmapFiles(m.fragQueue, true)
	return nil
}


func (m *BitmapIndex) BatchSetBit(stream pb.BitmapIndex_BatchSetBitServer) error {

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
		ts := time.Unix(0, kv.Time)

		select {
		case m.fragQueue <- newBitmapFragment(indexName, fieldName, rowID, ts, kv.Value):
		default:
			// Fragment queue is full, send to disk
			if err := m.journalBitmapFragment(kv.Value, indexName, fieldName, rowID, ts); err != nil {
				log.Printf("%s", err)
				return err
			}
		}
    }
}


func (m *BitmapIndex) journalBitmapFragment(newBm []byte, indexName, fieldName string, rowID uint64,
		ts time.Time) error {

	m.fragFileLock.Lock()
	defer m.fragFileLock.Unlock()

	if fd, err := m.newFragmentFile(indexName, fieldName, rowID, ts); err == nil {
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


func (m *BitmapIndex) saveCompleteBitmap(bm *StandardBitmap, indexName, fieldName string, rowID uint64, 
		ts time.Time, create bool) error {

	data, err := bm.Bits.MarshalBinary() 
	if err != nil {
		return err
	}

	if fd, err := m.writeCompleteFile(indexName, fieldName, rowID, ts, bm.TQType, create); err == nil {
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
	TQType		string
}


func newStandardBitmap(timeQuantumType string) *StandardBitmap {
	return &StandardBitmap{Bits: roaring.NewBitmap(), ModTime: time.Now().Add(time.Second * 10),
		TQType: timeQuantumType}
}


type BitmapFragment struct {
	IndexName			string
	FieldName			string
	RowID				uint64
	Time				time.Time
	BitData				[]byte
	ModTime				time.Time
}


func newBitmapFragment(index, field string, rowID uint64, ts time.Time, f []byte) *BitmapFragment {
	return &BitmapFragment{IndexName: index, FieldName: field, RowID: rowID, Time: ts, BitData: f}
}


func (m *BitmapIndex) generateFilePath(index, field string, rowID uint64, ts time.Time, 
		tqType string) string {

	baseDir := m.dataDir + SEP + "bitmap" + SEP + index + SEP + field + SEP + fmt.Sprintf("%d", rowID)
	fname := "default"
	if tqType == "YMD" {
		fname = ts.Format("2006-01-02T15")
	}
	if tqType == "YMDH" {
		baseDir = baseDir + SEP + fmt.Sprintf("%d%02d%02d", ts.Year(), ts.Month(), ts.Day())
		fname = ts.Format("2006-01-02T15")
	}
    os.MkdirAll(baseDir, 0755)
	return baseDir + SEP + fname
}


func (m *BitmapIndex) shouldWriteFile(index, field string, rowID uint64, ts time.Time, 
		tqType string, mod time.Time) (create bool, update bool) {

    path := m.generateFilePath(index, field, rowID, ts, tqType)
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


func (m *BitmapIndex) writeCompleteFile(index, field string, rowID uint64, ts time.Time, 
		tqType string, create bool) (*os.File, error) {

    path := m.generateFilePath(index, field, rowID, ts, tqType)
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


func (m *BitmapIndex) newFragmentFile(index, field string, rowID uint64, ts time.Time) (*os.File, error) {

	baseDir := m.dataDir + SEP + "overflow" + SEP + index + SEP + field + SEP +
		fmt.Sprintf("%d", rowID) + SEP + ts.Format("2006-01-02T15")
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


func (m *BitmapIndex) getTimeQuantum(index, field string) string {

	m.tableCacheLock.Lock()
	defer m.tableCacheLock.Unlock()

	var err error
	table, ok := m.tableCache[index]
	if !ok {
		configPath := m.EndPoint.dataDir + SEP + "config"
		//m.tableCacheLock.Lock()
		if table, err = shared.LoadSchema(configPath, index, "", m.EndPoint.consul); err != nil {
			log.Printf("ERROR: Could not load schema for %s - %v", index, err)
		} else {
			m.tableCache[index] = table
		}
		//m.tableCacheLock.Unlock()
	}
	timeQuantum := table.TimeQuantumType
	attr, err2 := table.GetAttribute(field)
	if err2 != nil {
		log.Printf("ERROR: Non existant attribute %s was referenced.", field)
	}
	if attr.TimeQuantumType != "" {
		timeQuantum = attr.TimeQuantumType
	}
	return timeQuantum
}


func (m *BitmapIndex) readBitmapFiles(fragQueue chan *BitmapFragment, isComplete bool) error {

	m.fragFileLock.Lock()
	defer m.fragFileLock.Unlock()

	var subPath string
    if isComplete {
		subPath = "bitmap"
	} else {
		subPath = "overflow"
	}

	baseDir := m.dataDir + SEP + subPath

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

		trPath := strings.Replace(path, baseDir + SEP, "", 1)

		s := strings.Split(trPath, SEP)
		if len(s) < 4 {
			err := fmt.Errorf("readBitmapFiles: Could not parse path [%s]", path)
			log.Println(err)
			return err
		}
		bf.IndexName = s[0]
		bf.FieldName = s[1]
		tq := m.getTimeQuantum(bf.IndexName, bf.FieldName)
		bf.RowID, err = strconv.ParseUint(s[2], 10, 64)
		if err != nil {
			err := fmt.Errorf("readBitmapFiles: Could not parse RowID - %v", err)
			log.Println(err)
			return err
		}

		if isComplete {
			if tq != "" {
				ts, err := time.Parse("2006-01-02T15", s[len(s) - 1])
				if err != nil {
					err := fmt.Errorf("readBitmapFiles: Could not parse '%s' Time[%s] - %v", s[len(s) - 1], tq, err)
					log.Println(err)
					return err
				}
				bf.Time = ts
			}
		} else {
			ts, err := time.Parse("2006-01-02T15", s[3])
			if err != nil {
				err := fmt.Errorf("readBitmapFiles: Could not parse '%s' Time - %v", s[len(s) - 1], err)
				log.Println(err)
				return err
			}
			bf.Time = ts
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


func (m *BitmapIndex) batchLoadProcessLoop(threadID int) {

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


func (m *BitmapIndex) overflowProcessLoop() {

	log.Printf("overflowProcessLoop - Started.")
	for {
		if err := m.readBitmapFiles(m.fragQueue, false); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second)
	}
	log.Printf("overflowProcessLoop - Ended.")
}


func (m *BitmapIndex) updateBitmapCache(f *BitmapFragment) {

	// If the rowID exists then merge in the new set of bits

   	start := time.Now()
	newBm := newStandardBitmap(m.getTimeQuantum(f.IndexName, f.FieldName))
	if err := newBm.Bits.UnmarshalBinary(f.BitData); err != nil {
		log.Printf("updateBitmapCache - UnmarshalBinary error - %v", err)
		return
	}
	m.bitmapCacheLock.Lock()
	if _, ok := m.bitmapCache[f.IndexName]; !ok {
		m.bitmapCache[f.IndexName] = make(map[string]map[uint64]map[time.Time]*StandardBitmap)
	}
	if _, ok := m.bitmapCache[f.IndexName][f.FieldName]; !ok {
		m.bitmapCache[f.IndexName][f.FieldName] = make(map[uint64]map[time.Time]*StandardBitmap)
	}
	if _, ok := m.bitmapCache[f.IndexName][f.FieldName][f.RowID]; !ok {
		m.bitmapCache[f.IndexName][f.FieldName][f.RowID] = make(map[time.Time]*StandardBitmap)
	}
	if existBm, ok := m.bitmapCache[f.IndexName][f.FieldName][f.RowID][f.Time]; !ok {
		m.bitmapCache[f.IndexName][f.FieldName][f.RowID][f.Time] = newBm
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


func (m *BitmapIndex) checkPersistBitmapCache() {

	m.bitmapCacheLock.RLock()
	defer m.bitmapCacheLock.RUnlock()

    for indexName, index := range m.bitmapCache {
        for fieldName, field := range index {
	        for rowID, ts := range field {
	            for t, bitmap := range ts {
					lastModSecs := time.Since(bitmap.ModTime).Round(time.Second)
					if dur, _ := time.ParseDuration("60s"); lastModSecs >= dur {
						bitmap.Lock.Lock()
	   					start := time.Now()
						if create, update := m.shouldWriteFile(indexName, fieldName, rowID, t, bitmap.TQType,
								bitmap.ModTime); create || update {
							if err := m.saveCompleteBitmap(bitmap, indexName, fieldName, rowID, t, create); err != nil {
								log.Printf("saveCompleteBitmap failed! - %v", err)
								bitmap.Lock.Unlock()
								return
							}
		    				elapsed := time.Since(start)
		    				log.Printf("Persist [%s/%s/%d/%s] done in %v. Last mod %v seconds ago.", indexName, 
							 	fieldName, rowID, t.Format("2006-01-02T15"), elapsed, lastModSecs)
						}
						bitmap.Lock.Unlock()
					}
				}
			}
		}
	}
}


func (m *BitmapIndex) Query(ctx context.Context, query *pb.BitmapQuery) (*wrappers.BytesValue, error) {

    if query == nil {
		return &wrappers.BytesValue{}, fmt.Errorf("query must not be nil")
	}
	
    if query.Query == nil {
		return &wrappers.BytesValue{}, fmt.Errorf("query fragment array must not be nil")
	}
    if len(query.Query) == 0 {
		return &wrappers.BytesValue{}, fmt.Errorf("query fragment array must not be empty")
	}
	fromTime := time.Unix(0, query.FromTime)
	toTime := time.Unix(0, query.ToTime)

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
		
		var bm *roaring.Bitmap
		var err error
		m.bitmapCacheLock.RLock()
		defer m.bitmapCacheLock.RUnlock()
		if bm, err = m.timeRange(v.Index, v.Field, v.RowID, fromTime, toTime); err != nil {
			return &wrappers.BytesValue{Value: []byte("")}, err
		}

		if firstTime  {
			prevOp = v.Operation
			result = bm
			firstTime = false
			continue
		}
		if prevOp == pb.QueryFragment_INTERSECT {
			result = roaring.ParAnd(0, result, bm)
		} else {
			result = roaring.ParOr(0, result, bm)
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


func (m *BitmapIndex) timeRange(index, field string, rowID uint64, fromTime, toTime time.Time) (*roaring.Bitmap, error) {
    
    tq := m.getTimeQuantum(index, field)
	result := roaring.NewBitmap()
	yr, mn, da := fromTime.Date()
	lookupTime := time.Date(yr, mn, da, 0, 0, 0, 0, time.UTC)

	if tq == "YMD" {
		for ; lookupTime.Before(toTime); lookupTime = lookupTime.AddDate(0, 0, 1) {
			if bm, ok := m.bitmapCache[index][field][rowID][lookupTime]; !ok {
				return nil, fmt.Errorf("Cannot find value for YMD [%s/%s/%d/%s]", index, field, rowID, 
						lookupTime.Format("2006-01-02T15"))
			} else {
				result = roaring.ParOr(0, result, bm.Bits)
			}
		}
	}
	if tq == "YMDH" {
		for ; lookupTime.Before(toTime); lookupTime = lookupTime.Add(time.Hour) {
			if bm, ok := m.bitmapCache[index][field][rowID][lookupTime]; !ok {
				return nil, fmt.Errorf("Cannot find value for YMDH [%s/%s/%d/%s]", index, field, rowID, 
						lookupTime.Format("2006-01-02T15"))
			} else {
				result = roaring.ParOr(0, result, bm.Bits)
			}
		}
	}
	return result, nil
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
