package generator

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/skycoin/services/otc/pkg/otc"
)

type Task func(*otc.User) (*otc.Order, error)

type Generator struct {
	UserCount int64
	Users     *sync.Map
	Logs      *log.Logger
	Task      Task
	Work      chan *otc.Work
}

func New(logs *log.Logger, task Task, work chan *otc.Work) *Generator {
	return &Generator{
		Users: &sync.Map{},
		Task:  task,
		Work:  work,
		Logs:  logs,
	}
}

func (g *Generator) Log(s string) {
	g.Logs.Println(s)
}

func (g *Generator) Count() int64 {
	return atomic.LoadInt64(&g.UserCount)
}

func (g *Generator) Tick() {
	g.Users.Range(g.Ranger(g.Task))
}

func (g *Generator) Add(user *otc.User) {
	_, exists := g.Users.LoadOrStore(user, nil)
	if !exists {
		// only add to count if didn't previously exist
		atomic.AddInt64(&g.UserCount, 1)
	}
}

func (g *Generator) Delete(user *otc.User) {
	atomic.AddInt64(&g.UserCount, -1)
	g.Users.Delete(user)
}

func (g *Generator) Ranger(task Task) func(k, v interface{}) bool {
	return func(k, v interface{}) bool {
		var user *otc.User = k.(*otc.User)

		// process user
		order, err := task(user)

		// log if error
		if err != nil {
			g.Logs.Println(err)
		}

		// if new order created, send to model
		if order != nil {
			// add to user
			user.Orders = append(user.Orders, order)

			// create work from order
			work := &otc.Work{
				Order: order,
				Done:  make(chan *otc.Result, 1),
			}

			// add err if any
			work.Return(err)

			// send to model
			g.Work <- work
		}

		// don't stop iterating over items
		return false
	}
}
