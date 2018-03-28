package model

import (
	"log"
	"os"

	"github.com/skycoin/services/otc/pkg/actor"
	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/monitor"
	"github.com/skycoin/services/otc/pkg/otc"
	"github.com/skycoin/services/otc/pkg/scanner"
	"github.com/skycoin/services/otc/pkg/sender"
	"github.com/skycoin/services/otc/pkg/watcher"
)

// TODO: use an interface for generators and actors
type Workers struct {
	Scanner *generator.Generator
	Sender  *actor.Actor
	Monitor *actor.Actor
}

func NewWorkers(curs *currencies.Currencies, watch *watcher.Watcher) (*Workers, chan *otc.Work) {
	work := make(chan *otc.Work, 0)

	return &Workers{
		Scanner: generator.New(
			log.New(os.Stdout, "[SCANNER] ", log.LstdFlags),
			scanner.Task(watch),
			work,
		),
		Sender: actor.New(
			log.New(os.Stdout, " [SENDER] ", log.LstdFlags),
			sender.Task(curs),
		),
		Monitor: actor.New(
			log.New(os.Stdout, "[MONITOR] ", log.LstdFlags),
			monitor.Task(curs),
		),
	}, work
}

func (w *Workers) Route(work *otc.Work) {
	switch work.Order.Status {
	case otc.DEPOSIT:
		w.Scanner.Add(work)
	case otc.SEND:
		w.Sender.Add(work)
	case otc.CONFIRM:
		w.Monitor.Add(work)
	}
}
