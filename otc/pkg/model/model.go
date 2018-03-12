package model

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/skycoin/services/otc/pkg/actor"
	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
)

var ErrReqMissing error = errors.New("request missing")

type Model struct {
	sync.RWMutex

	Running bool
	Workers *Workers
	Logger  *log.Logger
	Router  *actor.Actor
	Lookup  map[string]*otc.Request
	Stops   map[*actor.Actor]chan struct{}
	stop    chan struct{}
}

func New(curs *currencies.Currencies) (*Model, error) {
	workers := NewWorkers(curs)

	model := &Model{
		Running: true,
		Workers: workers,
		Logger:  log.New(os.Stdout, "    [OTC] ", log.LstdFlags),
		Router: actor.New(
			log.New(os.Stdout, "  [MODEL] ", log.LstdFlags),
			Task(workers),
		),
		Lookup: make(map[string]*otc.Request),
		Stops: map[*actor.Actor]chan struct{}{
			workers.Scanner: make(chan struct{}),
			workers.Sender:  make(chan struct{}),
			workers.Monitor: make(chan struct{}),
		},
		stop: make(chan struct{}),
	}

	// load all requests from disk
	reqs, err := Load()
	if err != nil {
		return nil, err
	}

	// add to model for processing
	for _, req := range reqs {
		if err = model.Load(req); err != nil {
			return nil, err
		}
	}

	defer model.Start()
	return model, nil
}

func (m *Model) Run(w time.Duration, s chan struct{}, a *actor.Actor) {
	for {
		<-time.After(w)

		select {
		case <-s:
			a.Logs.Println("stopping")
			return
		default:
			if !m.Paused() {
				a.Tick()
			}
		}
	}
}

func (m *Model) Start() {
	wait := time.Second * 5

	// start model routing actor
	go m.Run(wait, m.stop, m.Router)

	// start service actors on tick loop
	go m.Run(wait, m.Stops[m.Workers.Scanner], m.Workers.Scanner)
	go m.Run(wait, m.Stops[m.Workers.Sender], m.Workers.Sender)
	go m.Run(wait, m.Stops[m.Workers.Monitor], m.Workers.Monitor)

	go func() {
		for {
			<-time.After(wait)

			m.Logger.Printf(
				`[%d] [%d] [%d]`,
				m.Workers.Scanner.Count(),
				m.Workers.Sender.Count(),
				m.Workers.Monitor.Count(),
			)
		}
	}()
}

func (m *Model) Stop() {
	// stop model routing
	m.stop <- struct{}{}

	// stop service actors
	for _, stop := range m.Stops {
		stop <- struct{}{}
	}
}

func (m *Model) Status(iden string) (otc.Status, int64, error) {
	m.RLock()
	defer m.RUnlock()
	if m.Lookup[iden] == nil {
		return "", 0, ErrReqMissing
	}

	m.Lookup[iden].Lock()
	defer m.Lookup[iden].Unlock()
	return m.Lookup[iden].Status, m.Lookup[iden].Times.UpdatedAt, nil
}

func (m *Model) Load(req *otc.Request) error {
	m.Lock()
	defer m.Unlock()
	m.Lookup[req.Iden()] = req

	req.Lock()
	defer req.Unlock()

	work := &otc.Work{
		Request: req,
		Done:    make(chan *otc.Result, 1),
	}

	m.Workers.Route(work)
	m.Router.Add(work)

	return nil
}

func (m *Model) Add(req *otc.Request) {
	m.Lock()
	defer m.Unlock()
	m.Lookup[req.Iden()] = req

	req.Lock()
	defer req.Unlock()

	res := &otc.Result{
		Finished: time.Now().UTC().Unix(),
		Err:      nil,
	}
	if err := Save(req, res); err != nil {
		m.Logger.Println(err)
	}
	if req.Status == otc.NEW {
		req.Status = otc.DEPOSIT
	}

	work := &otc.Work{
		Request: req,
		Done:    make(chan *otc.Result, 1),
	}
	work.Done <- res

	m.Router.Add(work)
}

func (m *Model) Paused() bool {
	m.RLock()
	defer m.RUnlock()
	return !m.Running
}

func (m *Model) Pause() {
	m.Lock()
	defer m.Unlock()
	m.Logger.Println("paused")
	m.Running = false
}

func (m *Model) Unpause() {
	m.Lock()
	defer m.Unlock()
	m.Logger.Println("running")
	m.Running = true
}

func (m *Model) Reqs() []otc.Request {
	m.RLock()
	defer m.RUnlock()

	reqs := make([]otc.Request, 0)
	for _, req := range m.Lookup {
		reqs = append(reqs, *req)
	}
	return reqs
}
