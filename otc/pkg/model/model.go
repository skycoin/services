package model

import (
	"log"
	"os"
	"time"

	"github.com/skycoin/services/otc/pkg/actor"
	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
	"github.com/skycoin/services/otc/pkg/watcher"
)

type Config struct {
	Currencies *currencies.Currencies
	Watcher    *watcher.Watcher
}

type Model struct {
	Controller *Controller
	Lookup     *Lookup
	Workers    *Workers
	Router     *actor.Actor
	Work       chan *otc.Work
	Logs       *log.Logger
}

func New(conf *Config) (*Model, error) {
	workers, work := NewWorkers(conf)
	stoppers := make([]chan struct{}, 4, 4)

	model := &Model{
		Controller: NewController(stoppers),
		Lookup:     NewLookup(),
		Workers:    workers,
		Router: actor.New(
			log.New(os.Stdout, "  [MODEL] ", log.LstdFlags),
			Task(workers),
		),
		Work: work,
		Logs: log.New(os.Stdout, "    [OTC] ", log.LstdFlags),
	}

	// load all users from disk
	users, err := Load()
	if err != nil {
		return nil, err
	}

	// add each user to model
	for _, user := range users {
		if err = model.Add(user); err != nil {
			return nil, err
		}
	}

	defer model.Start()
	return model, nil
}

func (m *Model) Run(d time.Duration, s chan struct{}, w Worker) {
	for {
		<-time.After(d)

		select {
		case <-s:
			w.Log("stopping")
			return
		default:
			if !m.Controller.Paused() {
				w.Tick()
			}
		}
	}
}

func (m *Model) Start() {
	wait := time.Second * 5

	// TODO: move to own function somewhere, add stopper
	//
	// this receives work from generator and adds to the model where it is
	// saved and routed accordingly
	go func() {
		for {
			m.Router.Add(<-m.Work)
		}
	}()

	// start model routing actor
	go m.Run(wait, m.Controller.Stoppers[0], m.Router)

	// start actors
	go m.Run(wait, m.Controller.Stoppers[2], m.Workers.Sender)
	go m.Run(wait, m.Controller.Stoppers[3], m.Workers.Monitor)

	// start order generator
	go m.Run(wait, m.Controller.Stoppers[1], m.Workers.Scanner)

	// logging
	go func() {
		for {
			<-time.After(wait)

			m.Logs.Printf(
				`[%d] [%d] [%d]`,
				m.Workers.Scanner.Count(),
				m.Workers.Sender.Count(),
				m.Workers.Monitor.Count(),
			)
		}
	}()
}

func (m *Model) Add(user *otc.User) error {
	// add user to lookup map for later access
	m.Lookup.AddUser(user)
	m.Lookup.AddStatus(user)

	// save user to disk
	if err := SaveUser(user); err != nil {
		return err
	}

	// route existing orders
	for _, order := range user.Orders {
		result := &otc.Result{time.Now().UTC().Unix(), nil}

		// save to disk
		if err := SaveOrder(order, result); err != nil {
			return err
		}

		// create work
		work := &otc.Work{order, make(chan *otc.Result, 1)}
		work.Done <- result

		// route work
		m.Router.Add(work)
	}

	// add user to generator to watch for new orders
	m.Workers.Scanner.Add(user)

	return nil
}

func (m *Model) Orders() []otc.Order {
	orders := m.Lookup.GetOrders()

	safe := make([]otc.Order, len(orders), len(orders))
	for i := range orders {
		safe[i] = *orders[i]
	}

	return safe
}

func (m *Model) Users() []otc.User {
	users := m.Lookup.GetUsers()

	safe := make([]otc.User, len(users), len(users))
	for i := range users {
		safe[i] = *users[i]
	}

	return safe
}
