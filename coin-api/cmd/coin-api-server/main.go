package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

var (
	srvaddr = flag.String("srv", "localhost:12345", "RPC listener address")
)

func init() {
	flag.Parse()
}

func main() {
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", *srvaddr)
	if e != nil {
		log.Fatalf("Couldn't start listening on port %s. Error %s", e.Error(), *srvaddr)
	}
	log.Println("Serving RPC handler")
	// TODO(stgleb): Add request timeouts for server
	server := &http.Server{}
	err := server.Serve(l)

	if err != nil {
		log.Fatalf("Error serving: %s", err)
	}
}
