package signal

import (
	"sync"

	"github.com/skycoin/net/factory"
	"github.com/skycoin/services/signal/op2c"
	"github.com/skycoin/services/signal/op2s"
)

var (
	DefaultServer = NewServer()
)

func Listen(address string) error {
	return DefaultServer.Listen(address)
}

func GetClient(id uint) (client *Client, ok bool) {
	return DefaultServer.GetClient(id)
}

type Server struct {
	factory *factory.TCPFactory

	clients     map[uint]*Client
	fieldsMutex sync.RWMutex
}

func NewServer() *Server {
	f := factory.NewTCPFactory()
	s := &Server{factory: f, clients: make(map[uint]*Client)}
	s.factory.AcceptedCallback = s.accept
	return s
}

func (s *Server) Listen(address string) error {
	return s.factory.Listen(address)
}

func (s *Server) accept(conn *factory.Connection) {
	var err error
	client := newClient(conn, op2s.OPS, op2c.RESPS)
	client.initRespChan()
	defer func() {
		if e := recover(); e != nil {
			conn.GetContextLogger().Errorf("accept recover err %v", e)
		}
		if err != nil {
			conn.GetContextLogger().Errorf("accept err %v", err)
		}
		client.Close()
		s.removeClient(client)
	}()
	go func() {
		s.addClient(client)
	}()
	for {
		select {
		case m, ok := <-conn.GetChanIn():
			if !ok {
				return
			}
			err = client.Operate(client, m)
			if err != nil {
				conn.GetContextLogger().Errorf("execute err: %v", err)
			}
		}
	}
}

func (s *Server) addClient(c *Client) {
	reg := c.GetReg()
	if reg == nil {
		return
	}
	id := reg.Id
	s.fieldsMutex.Lock()
	s.clients[id] = c
	s.fieldsMutex.Unlock()
}

func (s *Server) GetClient(id uint) (client *Client, ok bool) {
	s.fieldsMutex.RLock()
	defer s.fieldsMutex.RUnlock()
	client, ok = s.clients[id]
	return
}

func (s *Server) removeClient(c *Client) {
	reg := c.GetReg()
	if reg == nil {
		return
	}
	id := reg.Id
	s.fieldsMutex.Lock()
	defer s.fieldsMutex.Unlock()
	client, ok := s.clients[id]
	if !ok || client != c {
		return
	}
	delete(s.clients, id)
}
