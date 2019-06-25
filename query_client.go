package main

import (
	"io/ioutil"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/go-yaml/yaml"
    "github.com/guymolinari/dstash/client"
	pb "github.com/guymolinari/dstash/grpc"
)

type QueryFragment struct {
    Index        string    	`yaml:"index"`
    Field        string     `yaml:"field"`
    RowID        uint64     `yaml:"rowID"`
    Operation    string		`yaml:"op"`
}


type BitmapQuery struct {
    FromTime     string `yaml:"fromTime"`
    ToTime       string `yaml:"toTime"`
	Query        []*QueryFragment `yaml:"query,omitempty"`
}


func main() {

/*
	f1 := &QueryFragment{Index: "user360", Field: "is_registered", RowID: 0, Operation: "UNION"}
	f2 := &QueryFragment{Index: "user360", Field: "is_registered", RowID: 1, Operation: "UNION"}
	f3 := &QueryFragment{Index: "user360", Field: "gender", RowID: 1, Operation: "INTERSECT"}
	q := &BitmapQuery{Query: []*QueryFragment{f1, f2, f3}}
*/

	if len(os.Args) < 2 {
		log.Fatal("Must specify query .yaml file.")
	}

    conn := dstash.NewDefaultConnection()
    conn.ServicePort = 4000
    if err := conn.Connect(); err != nil {
        log.Fatal(err)
    }

	client := dstash.NewBitmapIndex(conn, 2000000)
  
    q := BitmapQuery{}

	data, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        fmt.Printf("error: %v", err)
    }
    
    err = yaml.Unmarshal(data, &q)
    if err != nil {
        fmt.Printf("error: %v", err)
    }

    d, err := yaml.Marshal(&q)
    if err != nil {
        fmt.Printf("error: %v", err)
    }
    fmt.Printf("--- t dump:\n%s\n\n", string(d))

	start := time.Now()
	count, err := client.CountQuery(toProto(&q))
	if err != nil {
		log.Fatalf("Error during server call to CountQuery - %v", err)
	}
	
	elapsed := time.Since(start)
	log.Printf("Results = %d, elapsed time %v", count, elapsed)
}


func toProto(query *BitmapQuery) *pb.BitmapQuery {

	frags := make([]*pb.QueryFragment, len(query.Query))
	for i, f := range query.Query {
		var op pb.QueryFragment_OpType
		if v, ok := pb.QueryFragment_OpType_value[f.Operation]; ok {
			switch v {
			case 0:
				op = pb.QueryFragment_INTERSECT
			case 1:
				op = pb.QueryFragment_UNION
			case 2:
				op = pb.QueryFragment_DIFFERENCE
			}
		}
		frags[i] = &pb.QueryFragment{Index: f.Index, Field: f.Field, RowID: f.RowID, Operation: op}
	}
	bq :=  &pb.BitmapQuery{Query: frags}
	if fromTime, err := time.Parse("2006-01-02T15", query.FromTime); err == nil {
		bq.FromTime = fromTime.UnixNano()
	}
	if toTime, err := time.Parse("2006-01-02T15", query.ToTime); err == nil {
		bq.ToTime = toTime.UnixNano()
	}
	return bq
}

