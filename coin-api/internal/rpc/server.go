package rpc

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	shutDownTimeout = time.Second * 1
	readTimeout     = time.Second * 60
	writeTimeout    = time.Second * 120
	idleTimeout     = time.Second * 600
)

// Server represents an RPC server for coin-api
type Server struct {
	Handlers map[string]func(req Request) *Response
	server   *http.Server
}

func NewServer(srvAddr string, handlers map[string]func(r Request) *Response) *Server {
	server := &http.Server{
		Addr:         srvAddr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
	s := &Server{}

	s.server = server
	s.server.Handler = s
	s.Handlers = handlers

	return s
}

// Start starts an coin-api RPC Server and runs until shutdown channel is closed
func (s *Server) Start() {
	l, e := net.Listen("tcp", s.server.Addr)
	if e != nil {
		log.Fatalf("Couldn't start listening on port %s. Error %s", e.Error(), s.server.Addr)
	}

	log.Printf("Starting server %s\n", s.server.Addr)
	if err := s.server.Serve(l); err != http.ErrServerClosed {
		log.Printf("Error while running RPC server: %s\n", err)
	}
}

func (s *Server) ShutDown() {
	log.Println("ShutDown server...")
	ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Error while closing server: %s\n", err)
	}

	log.Println("Stopped server.")
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// reading request body
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer func() {
		if err = r.Body.Close(); err != nil {
			panic(err)
		}
	}()

	var req Request
	if err = json.Unmarshal(reqBytes, &req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	endpoint := strings.TrimPrefix(r.URL.Path, "/")

	handler, ok := s.Handlers[endpoint]
	if !ok {
		log.Printf("not found %s\n", endpoint)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := handler(req)

	if response == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response.Error != nil {
		log.Printf("error: %s\n", response.Error)
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if _, err = w.Write(responseBytes); err != nil {
		log.Printf("Failed to write responseBytes: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
