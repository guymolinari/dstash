package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
    "github.com/akrylysov/pogreb"
    "github.com/golang/protobuf/ptypes/empty"
    "github.com/golang/protobuf/ptypes/wrappers"
    pb "github.com/guymolinari/dstash/grpc"
)

type KVStore struct {
	*EndPoint
    storeCache         map[string]*pogreb.DB
}


func NewKVStore(endPoint *EndPoint) (*KVStore, error) {

	e := &KVStore{EndPoint: endPoint}
	e.storeCache = make(map[string]*pogreb.DB)
    pb.RegisterKVStoreServer(endPoint.server, e)
	return e, nil
}


func (m *KVStore) getStore(index string) (db *pogreb.DB, err error) {

	var ok bool
	if db, ok = m.storeCache[index]; ok {
		return
	}

    db, err = pogreb.Open(m.EndPoint.dataDir + SEP + index + SEP + "data.dat", nil)
    if err != nil {
		m.storeCache[index] = db
    }
    return
}


// Put - Insert a new key
func (m *KVStore) Put(ctx context.Context, kv *pb.IndexKVPair) (*empty.Empty, error) {

    if kv == nil {
		return &empty.Empty{}, fmt.Errorf("KV Pair must not be nil")
	}
    if kv.Key == nil || len(kv.Key) == 0 {
		return &empty.Empty{}, fmt.Errorf("Key must be specified.")
	}
    if kv.IndexPath == "" {
		return &empty.Empty{}, fmt.Errorf("Index must be specified.")
	}
	db, err := m.getStore(kv.IndexPath)
	if err != nil {
    	return &empty.Empty{}, err
	}
	err = db.Put(kv.Key, kv.Value)
	if err != nil {
    	return &empty.Empty{}, err
	}
    return &empty.Empty{}, nil
}



func (m *KVStore) Lookup(ctx context.Context, kv *pb.IndexKVPair) (*pb.IndexKVPair, error) {
    if kv == nil {
		return &pb.IndexKVPair{}, fmt.Errorf("KV Pair must not be nil")
	}
    if kv.Key == nil || len(kv.Key) == 0 {
		return &pb.IndexKVPair{}, fmt.Errorf("Key must be specified.")
	}
    if kv.IndexPath == "" {
		return &pb.IndexKVPair{}, fmt.Errorf("Index must be specified.")
	}
	db, err := m.getStore(kv.IndexPath)
	if err != nil {
    	b := make([]byte, 8)
    	binary.LittleEndian.PutUint64(b, 0)
        kv.Value = b
		return kv, err
	}
	val, err := db.Get(kv.Key)
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


func (m *KVStore) BatchPut(stream pb.KVStore_BatchPutServer) error {

    var putCount int32
    for {
        kv, err := stream.Recv()
		db, err2 := m.getStore(kv.IndexPath)
		if err2 != nil {
			return err2
		}
        if err == io.EOF {
            db.Sync()
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
	    if kv.IndexPath == "" {
			return fmt.Errorf("Index must be specified.")
		}
		if err := db.Put(kv.Key, kv.Value); err != nil {
    		return err
		}
    }
}


func (m *KVStore) BatchLookup(stream pb.KVStore_BatchLookupServer) error {

    for {
        kv, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
		db, err := m.getStore(kv.IndexPath)
        if err != nil {
    		b := make([]byte, 8)
    		binary.LittleEndian.PutUint64(b, 0)
   	     	kv.Value = b
            return err
        }
		val, err := db.Get(kv.Key)
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


func (m *KVStore) Items(index *wrappers.StringValue, stream pb.KVStore_ItemsServer) error {

	if index.Value == "" {
		return fmt.Errorf("Index must be specified!")
	}
	db, err := m.getStore(index.Value)
	if err != nil {
		return err
	}

	it := db.Items()
	for {
	    key, val, err := it.Next()
	    if err != nil {
	        if err != pogreb.ErrIterationDone {
	            return err
	        }
	        break
	    }
        if err := stream.Send(&pb.IndexKVPair{IndexPath: index.Value, Key: key, Value: val}); err != nil {
            return err
        }
	}
    return nil
}
