package buyer

import (
	"container/list"
	"log"
	"os"
	"sync"
	"time"

	"github.com/skycoin/services/otc/types"
)

type Buyer struct {
	sync.Mutex

	work   *list.List
	stop   chan struct{}
	logger *log.Logger
}

func NewBuyer() (*Buyer, error) {
	return &Buyer{
		work:   list.New().Init(),
		stop:   make(chan struct{}),
		logger: log.New(os.Stdout, types.LOG_BUYER, types.LOG_FLAGS),
	}, nil
}

func (b *Buyer) Stop() { b.stop <- struct{}{} }

func (b *Buyer) Start() {
	b.logger.Println("started")

	go func() {
		for {
			// TODO: use config
			<-time.After(time.Second * 1)

			b.logger.Printf("[%d]\n", b.work.Len())

			select {
			case <-b.stop:
				b.logger.Println("stopped")
				return
			default:
				b.process()
			}
		}
	}()
}

func (b *Buyer) process() {
	b.Lock()
	defer b.Unlock()

	for e := b.work.Front(); e != nil; e = e.Next() {
		w := e.Value.(*types.Work)

		switch w.Request.Metadata.BuyStatus {
		case types.EXCHANGE_DEPOSIT:
			// get exchange deposit address and send
		case types.EXCHANGE_CONFIRM:
			// check exchange deposit address for balance
		case types.EXCHANGE_TRADE:
			// execute trade for request's amount
		case types.EXCHANGE_RETURN:
			// withdraw to otc
		case types.EXCHANGE_RETURNED:
			// check if withdraw was confirmed
			// move request to sender
		}
	}
}

func (b *Buyer) Handle(request *types.Request) chan *types.Result {
	b.Lock()
	defer b.Unlock()

	result := make(chan *types.Result, 1)
	b.work.PushBack(&types.Work{request, result})
	return result
}
