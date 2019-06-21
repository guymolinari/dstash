package main

import (
    "log"
    "os"
	"time"
    "gopkg.in/alecthomas/kingpin.v2"
	"github.com/guymolinari/dstash/server"
)

var (
    Version  string
    Build    string
)

func main() {
    app := kingpin.New(os.Args[0], "DStash server node.").DefaultEnvars()
    //app.Version("Version: " + Version + "\nBuild: " + Build)
    dataDir := app.Arg("data-dir", "Root directory for data files").Default("/home/ec2-user/data").String()
    bindAddr := app.Arg("bind", "Bind address for this endpoint.").Default("127.0.0.1").String()
    port := app.Arg("port", "Port for this endpoint.").Default("5000").Int32()
    tls := app.Flag("tls", "Connection uses TLS if true.").Bool()
    certFile := app.Flag("cert-file", "TLS cert file path.").String()
    keyFile := app.Flag("key-file", "TLS key file path.").String()

    kingpin.MustParse(app.Parse(os.Args[1:]))

    m, err:= server.NewEndPoint(*dataDir)
    if err != nil {
        log.Printf("[node: Cannot initialize endpoint config: error: %s", err)
    }
    m.BindAddr = *bindAddr
    m.Port = uint(*port)
    _ = *tls
    _ = *certFile
    _ = *keyFile

	_, err2 := server.NewKVStore(m)
    if err2 != nil {
        log.Printf("[node: Cannot initialize kv store config: error: %s", err2)
    }

	_, err3 := server.NewStringSearch(m)
    if err3 != nil {
        log.Printf("[node: Cannot initialize search config: error: %s", err3)
    }

	start := time.Now()
	bitmapIndex := server.NewBitmapIndex(m)
	bitmapIndex.Init()
	elapsed := time.Since(start)
	log.Printf("Bitmap index initialized in %v.", elapsed)

    node, err := server.Join("dstash", m)
    if err != nil {
    	log.Printf("[node: Cannot initialize endpoint config: error: %s", err)
	}
    <-node.Stop
    err = <- node.Err
    if err != nil {
    	log.Printf("[node: Cannot initialize endpoint config: error: %s", err)
	}
}

