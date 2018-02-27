package model

import (
	"container/list"
	"sync"
	"testing"
	"time"

	"github.com/skycoin/services/otc/types"
)

var (
	SCANNER types.Service
	SENDER  types.Service
	MONITOR types.Service
)

type Service struct {
	sync.Mutex

	next types.Status
	work *list.List
	stop chan struct{}
}

func (s *Service) Stop() { s.stop <- struct{}{} }

func (s *Service) Start() {
	go func() {
		for {
			<-time.After(time.Second / 4)

			select {
			case <-s.stop:
				return
			default:
				s.process()
			}
		}
	}()
}

func (s *Service) process() {
	s.Lock()
	defer s.Unlock()

	for e := s.work.Front(); e != nil; e = e.Next() {
		w := e.Value.(*types.Work)

		// just send to next service for testing
		w.Request.Metadata.Status = s.next
		w.Return(nil)
		s.work.Remove(e)
	}
}

func (s *Service) Handle(r *types.Request) chan *types.Result {
	s.Lock()
	defer s.Unlock()

	result := make(chan *types.Result, 1)
	s.work.PushBack(&types.Work{r, result})
	return result
}

func NewService(s types.Status) types.Service {
	return &Service{
		next: s,
		work: list.New().Init(),
		stop: make(chan struct{}),
	}
}

func init() {
	SCANNER = NewService(types.SEND)
	SENDER = NewService(types.CONFIRM)
	MONITOR = NewService(types.DONE)
}

func TestNewModel(t *testing.T) {
	// good
	_, err := NewModel(
		&types.Config{
			Model: struct {
				Tick int
				Path string
			}{
				Tick: 1,
				Path: "../db/",
			},
		},
		SCANNER,
		SENDER,
		MONITOR,
	)
	if err != nil {
		t.Fatal(err)
	}

	// nil service
	_, err = NewModel(
		&types.Config{
			Model: struct {
				Tick int
				Path string
			}{
				Tick: 1,
				Path: "../db/",
			},
		},
		nil,
		SENDER,
		MONITOR,
	)
	if err != ErrNilService {
		t.Fatal("didn't handle nil service")
	}

	// bad path
	_, err = NewModel(
		&types.Config{
			Model: struct {
				Tick int
				Path string
			}{
				Tick: 1,
				Path: "",
			},
		},
		SCANNER,
		SENDER,
		MONITOR,
	)
	if err == nil {
		t.Fatal("didn't handle bad path")
	}
}

func TestModelAdd(t *testing.T) {
	m, err := NewModel(
		&types.Config{
			Model: struct {
				Tick int
				Path string
			}{
				Tick: 1,
				Path: "testing/",
			},
		},
		SCANNER,
		SENDER,
		MONITOR,
	)
	if err != nil {
		t.Fatal(err)
	}

	m.Add(&types.Request{
		Address:  types.Address("2dvVgeKNU7UHdvvBUVZXbBaxoTkpemo1cmg"),
		Currency: types.Currency("testing"),
		Drop:     types.Drop("testing"),
		Metadata: &types.Metadata{
			Status:    types.DEPOSIT,
			CreatedAt: 0,
			UpdatedAt: 0,
			TxId:      "",
		},
	})

	if m.results.Len() != 1 {
		t.Fatal("model not queueing work properly")
	}

	// check on disk
}
