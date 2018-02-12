package main

import (
	"flag"
	"github.com/skycoin/services/coin-api/rpc"
	"log"
	"net"
	"net/http"
)

var (
	srvaddr = flag.String("srv", "localhost:12345", "RPC listener address")
)

func init() {
	flag.Parse()
}

func main() {

	l, e := net.Listen("tcp", *srvaddr)
	if e != nil {
		log.Fatalf("Couldn't start listening on port %s. Error %s", e.Error(), *srvaddr)
	}
	log.Println("Serving RPC handler")
	// TODO(stgleb): Add request timeouts for server

	server := &http.Server{
		Addr: *srvaddr,
	}
	err := server.Serve(l)
	shutDownChan := make(chan struct{})
	rpcServer := rpc.NewServer(server, map[string]func(request rpc.Request) *rpc.Response{}, shutDownChan)

	rpcServer.Start()

	if err != nil {
		log.Fatalf("Error serving: %s", err)
	}
}
