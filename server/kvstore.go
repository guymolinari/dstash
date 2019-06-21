package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
    "github.com/akrylysov/pogreb"
    "github.com/golang/protobuf/ptypes/empty"
    pb "github.com/guymolinari/dstash/grpc"
)

type KVStore struct {
	*EndPoint
    store         *pogreb.DB
}


func NewKVStore(endPoint *EndPoint) (*KVStore, error) {

    db, err := pogreb.Open(endPoint.dataDir + "/" + "default.dat", nil)
    if err != nil {
        return nil, err
    }

	e := &KVStore{EndPoint: endPoint, store: db}
    pb.RegisterKVStoreServer(endPoint.server, e)
	return e, nil
}


// Put - Insert a new key
func (m *KVStore) Put(ctx context.Context, kv *pb.KVPair) (*empty.Empty, error) {
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



func (m *KVStore) Lookup(ctx context.Context, kv *pb.KVPair) (*pb.KVPair, error) {
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


func (m *KVStore) BatchPut(stream pb.KVStore_BatchPutServer) error {

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


func (m *KVStore) BatchLookup(stream pb.KVStore_BatchLookupServer) error {

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


func (m *KVStore) Items(e *empty.Empty, stream pb.KVStore_ItemsServer) error {

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
