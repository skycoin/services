package sender

import (
	"container/list"
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/skycoin/services/otc/dropper"
	"github.com/skycoin/services/otc/skycoin"
	"github.com/skycoin/services/otc/types"
	"github.com/skycoin/skycoin/src/api/cli"
)

type Sender struct {
	sync.Mutex

	config         *types.Config
	skycoin        *skycoin.Connection
	dropper        *dropper.Dropper
	logger         *log.Logger
	work           *list.List
	stop           chan struct{}
	fromAddrs      []string
	fromChangeAddr string
}

var ErrZeroBalance = errors.New("sender got drop with zero balance")

func NewSender(c *types.Config, s *skycoin.Connection, d *dropper.Dropper) (*Sender, error) {
	sender := &Sender{
		config:  c,
		skycoin: s,
		dropper: d,
		logger:  log.New(os.Stdout, types.LOG_SENDER, types.LOG_FLAGS),
		work:    list.New().Init(),
		stop:    make(chan struct{}),
	}

	sender.fromAddrs = sender.getFromAddrs()
	sender.fromChangeAddr = sender.fromAddrs[0]

	return sender, nil
}

func (s *Sender) Stop() { s.stop <- struct{}{} }

func (s *Sender) Start() {
	s.logger.Println("started")

	go func() {
		for {
			<-time.After(time.Second * time.Duration(s.config.Sender.Tick))

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

func (s *Sender) process() {
	s.Lock()
	defer s.Unlock()

	for e := s.work.Front(); e != nil; e = e.Next() {
		// convert list element to work
		w := e.Value.(*types.Work)

		// get value of amount
		value, err := s.dropper.GetValue(
			w.Request.Currency,
			// filled by Scanner
			w.Request.Metadata.Amount,
		)
		if err != nil {
			w.Return(err)
			s.work.Remove(e)
			continue
		}

		// create sky transaction
		tx, err := cli.CreateRawTx(
			s.skycoin.Client,
			s.skycoin.Wallet,
			s.fromAddrs,
			s.fromChangeAddr,
			[]cli.SendAmount{{Addr: string(w.Request.Address), Coins: value}},
		)
		if err != nil {
			w.Return(err)
			s.work.Remove(e)
			continue
		}

		// inject and get txId
		txId, err := s.skycoin.Client.InjectTransaction(tx)
		if err != nil {
			w.Return(err)
			s.work.Remove(e)
			continue
		}

		// next step is monitor service
		w.Request.Metadata.TxId = txId
		w.Request.Metadata.Status = types.CONFIRM
		w.Return(nil)
		s.work.Remove(e)
	}
}

func (s *Sender) getFromAddrs() []string {
	// get all addrs from wallet
	addrs := s.skycoin.Wallet.GetAddresses()

	// convert to string slice
	out := make([]string, len(addrs), len(addrs))
	for i := range addrs {
		out[i] = addrs[i].String()
	}

	return out
}

func (s *Sender) Handle(request *types.Request) chan *types.Result {
	s.Lock()
	defer s.Unlock()

	result := make(chan *types.Result, 1)
	s.work.PushBack(&types.Work{request, result})
	return result
}
