package actor

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/skycoin/services/otc/pkg/otc"
)

type Task func(*otc.Work) (bool, error)

type Actor struct {
	Reqs int64
	Work *sync.Map
	Task Task
	Logs *log.Logger
}

func New(logs *log.Logger, task Task) *Actor {
	return &Actor{0, &sync.Map{}, task, logs}
}

func (a *Actor) Count() int64 { return atomic.LoadInt64(&a.Reqs) }
func (a *Actor) Tick()        { a.Work.Range(a.Ranger(a.Task)) }

func (a *Actor) Add(work *otc.Work) {
	if _, existed := a.Work.LoadOrStore(work, nil); !existed {
		atomic.AddInt64(&a.Reqs, 1)
	}
}

func (a *Actor) Delete(work *otc.Work) {
	atomic.AddInt64(&a.Reqs, -1)

	a.Work.Delete(work)
}

func (a *Actor) Ranger(task Task) func(k, v interface{}) bool {
	return func(k, v interface{}) bool {
		var work *otc.Work = k.(*otc.Work)

		// process work
		done, err := task(work)

		// log if error
		if err != nil {
			a.Logs.Println(err)
		}

		if done {
			// delete from actor map
			a.Delete(work)
			// return to model for saving
			work.Return(err)
		}

		// don't stop ranging
		return true
	}
}
