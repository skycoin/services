package main

import (
	"flag"
	"fmt"
	. "github.com/skycoin/services/coin-api"
	"os"
	"os/signal"
	"syscall"
)

var (
	srvaddr = flag.String("srv", "localhost:12345", "RPC listener address")
)

func init() {
	flag.Parse()
}

func main() {
	// Add handlers for all currencies here
	handlers := map[string]func(request Request) *Response{
		"btc": BtcHandler,
	}

	// Create new server
	rpcServer := NewServer(*srvaddr, handlers)
	// Register shutdown handler
	registerShutdownHandler(rpcServer)
	// Start server
	rpcServer.Start()
}

func registerShutdownHandler(server *Server) {
	go func() {
		interruptChannel := make(chan os.Signal, 1)
		signal.Notify(interruptChannel, syscall.SIGINT)

		// Listen for initial shutdown signal and close the returned
		// channel to notify the caller.
		select {
		case sig := <-interruptChannel:
			fmt.Printf("Received signal (%s).  Shutting down...\n", sig)
			server.ShutDown()
		}

		// Listen for repeated signals and display a message so the user
		// knows the shutdown is in progress and the process is not
		// hung.
		for {
			select {
			case sig := <-interruptChannel:
				fmt.Printf("Received signal (%s).  Already "+
					"shutting down...", sig)
				os.Exit(1)
			}
		}
	}()
}
