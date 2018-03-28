package model

import (
	"log"
	"os"

	"github.com/skycoin/services/otc/pkg/actor"
	"github.com/skycoin/services/otc/pkg/generator"
	"github.com/skycoin/services/otc/pkg/monitor"
	"github.com/skycoin/services/otc/pkg/otc"
	"github.com/skycoin/services/otc/pkg/scanner"
	"github.com/skycoin/services/otc/pkg/sender"
)

type Worker interface {
	Tick()
	Log(string)
}

type Workers struct {
	Scanner *generator.Generator
	Sender  *actor.Actor
	Monitor *actor.Actor
}

func NewWorkers(conf *Config) (*Workers, chan *otc.Work) {
	work := make(chan *otc.Work, 0)

	return &Workers{
		Scanner: generator.New(
			log.New(os.Stdout, "[SCANNER] ", log.LstdFlags),
			scanner.Task(conf.Watcher),
			work,
		),
		Sender: actor.New(
			log.New(os.Stdout, " [SENDER] ", log.LstdFlags),
			sender.Task(conf.Currencies),
		),
		Monitor: actor.New(
			log.New(os.Stdout, "[MONITOR] ", log.LstdFlags),
			monitor.Task(conf.Currencies),
		),
	}, work
}

func (w *Workers) Route(work *otc.Work) {
	switch work.Order.Status {
	case otc.SEND:
		w.Sender.Add(work)
	case otc.CONFIRM:
		w.Monitor.Add(work)
	}
}
