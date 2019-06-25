package server

import (
    "context"
    "fmt"
    "log"
    "net"
    "os"
    "path"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/testdata"
	"github.com/hashicorp/consul/api"
    "github.com/golang/protobuf/ptypes/empty"
    "google.golang.org/grpc/health/grpc_health_v1"
    pb "github.com/guymolinari/dstash/grpc"
)

const (
	SEP = string(os.PathSeparator)
	
)

type EndPoint struct {
    Port 		  uint
    BindAddr      string
    tls           bool
    certFile      string
    keyFile       string
	dataDir		  string
    hashKey       string
	server		  *grpc.Server
	consul    	  *api.Client
}


func NewEndPoint(dataDir string) (*EndPoint, error) {

    m := &EndPoint{}
    m.hashKey = path.Base(dataDir)  // leaf directory name is consisten hash key
    if m.hashKey == "" || m.hashKey == "/" {
        return nil, fmt.Errorf("Data dir must not be root.")
	}

	m.dataDir = dataDir

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
            return nil, fmt.Errorf("Failed to generate credentials %v", err)
        }
        opts = append(opts, grpc.Creds(creds))
    }

    m.server = grpc.NewServer(opts...)
    pb.RegisterClusterAdminServer(m.server, m)
    grpc_health_v1.RegisterHealthServer(m.server, &HealthImpl{})

    return m, nil
}


func (m *EndPoint) Start() error {

    lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", m.BindAddr, m.Port))
    if err != nil {
        return fmt.Errorf("failed to listen: %v", err)
    }

    m.server.Serve(lis)
    return nil
}


func (m *EndPoint) Status(ctx context.Context, e *empty.Empty) (*pb.StatusMessage, error) {
    log.Printf("Status ping returning OK.\n")
    return &pb.StatusMessage{Status: "OK"}, nil
}


