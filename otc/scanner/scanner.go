package scanner

import (
	"container/list"
	"log"
	"os"
	"sync"
	"time"

	"github.com/skycoin/services/otc/dropper"
	"github.com/skycoin/services/otc/types"
)

type Scanner struct {
	sync.Mutex

	config  *types.Config
	dropper *dropper.Dropper
	logger  *log.Logger
	errs    *log.Logger
	work    *list.List
	stop    chan struct{}
}

func NewScanner(conf *types.Config, drpr *dropper.Dropper) (*Scanner, error) {
	return &Scanner{
		dropper: drpr,
		logger:  log.New(os.Stdout, types.LOG_SCANNER, types.LOG_FLAGS),
		errs:    log.New(os.Stdout, types.LOG_ERRS, types.LOG_FLAGS),
		config:  conf,
		work:    list.New().Init(),
		stop:    make(chan struct{}),
	}, nil
}

func (s *Scanner) Stop() { s.stop <- struct{}{} }

func (s *Scanner) Start() {
	s.logger.Println("started")

	go func() {
		for {
			<-time.After(time.Second * time.Duration(s.config.Scanner.Tick))

			s.logger.Printf("[%d]\n", s.work.Len())

			select {
			case <-s.stop:
				s.logger.Println("stopped")
				return
			default:
				s.process()
			}
		}
	}()
}

func (s *Scanner) process() {
	s.Lock()
	defer s.Unlock()

	for e := s.work.Front(); e != nil; e = e.Next() {
		w := e.Value.(*types.Work)

		// get balance of drop
		bal, err := s.dropper.GetBalance(w.Request.Currency, w.Request.Drop)
		if err != nil {
			w.Return(err)
			s.work.Remove(e)
			continue
		}

		// if user made a deposit
		if bal != 0 {
			w.Request.Metadata.Status = types.SEND
			w.Request.Metadata.Amount = bal
			w.Return(nil)
			s.work.Remove(e)
			continue
		}

		// balance is 0, so check if expired
		if w.Request.Metadata.Expired(s.config.Scanner.Expiration) {
			w.Request.Metadata.Status = types.EXPIRED
			w.Return(nil)
			s.work.Remove(e)
			continue
		}
	}
}

func (s *Scanner) Handle(request *types.Request) chan *types.Result {
	s.Lock()
	defer s.Unlock()

	result := make(chan *types.Result, 1)
	s.work.PushBack(&types.Work{request, result})
	return result
}
