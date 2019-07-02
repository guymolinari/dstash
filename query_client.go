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
    Operation    string		`yaml:"op,omitempty"`
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

	var usage error
	if len(os.Args) < 3 {
		usage = fmt.Errorf("Must specify query .yaml file.")
	}
	var action string
	switch os.Args[1] {
	case "count":
		action = "count"
	case "results":
		action = "results"
	default:
		usage = fmt.Errorf("action must be 'count' or 'results'.")
	}
	if usage != nil {
		log.Printf("%v", usage)
		log.Fatalf("Usage query_client count|results [yaml_file]")
	}

    conn := dstash.NewDefaultConnection()
    conn.ServicePort = 4000
    if err := conn.Connect(); err != nil {
        log.Fatal(err)
    }

	client := dstash.NewBitmapIndex(conn, 2000000)
  
    q := BitmapQuery{}

	data, err := ioutil.ReadFile(os.Args[2])
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
    fmt.Printf("vvv query dump:\n%s\n\n", string(d))

	if action == "count" {
		start := time.Now()
		count, err := client.CountQuery(toProto(&q))
		if err != nil {
			log.Fatalf("Error during server call to CountQuery - %v", err)
		}
		elapsed := time.Since(start)
		log.Printf("Results = %d, elapsed time %v", count, elapsed)
	}
	if action == "results" {
		start := time.Now()
		results, err := client.ResultsQuery(toProto(&q), 10000)
		if err != nil {
			log.Fatalf("Error during server call to ResultsQuery - %v", err)
		}
		elapsed := time.Since(start)
		for i := 0; i < 10000 && i < len(results); i++ {
			fmt.Printf("%d, ", results[i])
		}
		log.Println("Done.")
		log.Printf("Elapsed time %v", elapsed)
	}
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

