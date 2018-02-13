package rpc

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// Server represents an RPC server for coin-api
type Server struct {
	Handlers map[string]func(req Request) *Response
	TimeOut  time.Duration
	server   *http.Server
	shutDown chan struct{}
}

func NewServer(server *http.Server, handlers map[string]func(r Request) *Response, shutDown chan struct{}) *Server {
	s := &Server{}
	var mux = http.NewServeMux()

	// Register all handlers
	for k := range handlers {
		mux.HandleFunc("/"+k, s.handler)
	}

	server.Handler = mux
	s.server = server
	s.shutDown = shutDown

	return s
}

// Start starts an coin-api RPC Server and runs until shutdown channel is closed
func (s *Server) Start() {
	go func() {
		log.Printf("Starting server %s\n", s.server.Addr)
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Error while running RPC server: %s\n", err)
		}
	}()

	go func() {
		<-s.shutDown

		log.Println("ShutDown server...")
		ctx, cancel := context.WithTimeout(context.Background(), s.TimeOut)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			log.Printf("Error while closing server: %s\n", err)
		}

		log.Println("Stopped server.")
	}()
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {

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
