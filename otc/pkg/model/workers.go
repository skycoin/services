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
)

type Workers struct {
	Scanner *actor.Actor
	Sender  *actor.Actor
	Monitor *actor.Actor
}

func NewWorkers(curs *currencies.Currencies) *Workers {
	return &Workers{
		Scanner: actor.New(
			log.New(os.Stdout, "[SCANNER] ", log.LstdFlags),
			scanner.Task(curs),
		),
		Sender: actor.New(
			log.New(os.Stdout, " [SENDER] ", log.LstdFlags),
			sender.Task(curs),
		),
		Monitor: actor.New(
			log.New(os.Stdout, "[MONITOR] ", log.LstdFlags),
			monitor.Task(curs),
		),
	}
}

func (w *Workers) Route(work *otc.Work) {
	switch work.Request.Status {
	case otc.DEPOSIT:
		w.Scanner.Add(work)
	case otc.SEND:
		w.Sender.Add(work)
	case otc.CONFIRM:
		w.Monitor.Add(work)
	}
}
