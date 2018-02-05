package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/skycoin/skycoin/src/api/cli"
	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

type StepType int

const (
	// MERGE decribes a transaction from many addresses to 1.
	MERGE StepType = iota
	// SPLIT describes a transaction from 1 address to many.
	SPLIT
	// CLEAN describes a transaction from all addresses to 1.
	CLEAN
)

const (
	// SIZE is the number of addresses to SPLIT across. If SIZE is 2, then steps
	// will alternate between 1->2 and 2->1. If SIZE is 10, then steps will
	// alternate between 1->10 and 10->1.
	SIZE = 2
)

var (
	ErrNotEnoughAddrs = errors.New("not enough addresses to split across")
	ErrNoBalance      = errors.New("at least one address must have 1 SKY")
	ErrSmallBalance   = errors.New("at least one address must have 1 SKY")
)

type StepResult struct {
	Id       string
	Type     StepType
	From     []string
	To       []cli.SendAmount
	Start    time.Time
	End      time.Time
	Duration time.Duration
	Status   *visor.TransactionStatus
}

func (r *StepResult) String() string {
	var out string

	if r.Type == SPLIT {
		out = out + "SPLIT:\n"
	} else if r.Type == MERGE {
		out = out + "MERGE:\n"
	} else if r.Type == CLEAN {
		out = out + "CLEAN:\n"
	}

	var from string
	for i := range r.From {
		if i != 0 {
			from = from + "            "
		}
		from = from + r.From[i] + "\n"
	}
	out = out + "- From:     " + from

	var to string
	for i := range r.To {
		if i != 0 {
			to = to + "            "
		}
		to = to + fmt.Sprintf(
			"%s\t(%0.2fSKY)\n",
			r.To[i].Addr,
			float32(float32(r.To[i].Coins)/1e6),
		)
	}
	out = out + "- To:       " + to

	out = out + fmt.Sprintf("- Start:    %v\n", r.Start)
	out = out + fmt.Sprintf("- End:      %v\n", r.End)
	out = out + fmt.Sprintf("- Duration: %v\n", r.Duration)

	return out
}

type Step struct {
	Client *webrpc.Client
	Wallet *wallet.Wallet
	Logger *log.Logger

	Type   StepType
	Addrs  []string
	Amount uint64
	// To contains the addresses that will receive coins for the current step.
	To []string
	// From contains the addresses that will send coins for the current step.
	From []string
}

func NewStep(c *webrpc.Client, w *wallet.Wallet, l *log.Logger,
	addrs []cipher.Address) (*Step, error) {
	if len(addrs) < SIZE {
		return nil, ErrNotEnoughAddrs
	}

	s := &Step{
		Client: c,
		Wallet: w,
		Logger: l,
		Addrs:  nil,
		Type:   MERGE,
	}

	// convert []cipher.Address to []string for easier handling
	s.Addrs = make([]string, len(addrs), len(addrs))
	for i, addr := range addrs {
		s.Addrs[i] = addr.String()
	}

	// get initial s.To (which will become s.From on first run)
	var err error
	if s.To, err = s.getOrigin(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Step) getOrigin() ([]string, error) {
	unspent, err := s.Client.GetUnspentOutputs(s.Addrs)
	if err != nil {
		return nil, err
	}

	spendable, err := unspent.Outputs.SpendableOutputs().ToUxArray()
	if err != nil {
		return nil, err
	} else if len(spendable) == 0 {
		return nil, ErrNoBalance
	}

	var (
		max    uint64
		origin string
	)

	for i := range spendable {
		if spendable[i].Body.Coins > max {
			max = spendable[i].Body.Coins
			if max >= 1e6 {
				origin = spendable[i].Body.Address.String()
				break
			}
		}
	}

	// if couldn't find address with > 1 SKY
	if max < 1e6 {
		return nil, ErrSmallBalance
	}

	// find the id of initial address in the slice
	for _, addr := range s.Addrs {
		if addr == origin {
			return []string{addr}, nil
		}
	}

	return nil, ErrNoBalance
}

func (s *Step) getAmount(addrs []string) (uint64, error) {
	unspent, err := s.Client.GetUnspentOutputs(addrs)
	if err != nil {
		return 0, err
	}

	spendable, err := unspent.Outputs.SpendableOutputs().ToUxArray()
	if err != nil {
		return 0, err
	} else if len(spendable) == 0 {
		return 0, ErrNoBalance
	}

	var balance uint64

	for _, a := range spendable {
		balance = balance + a.Body.Coins
	}

	return round(balance), nil
}

func round(n uint64) uint64 {
	return ((n / 1e5) * 1e5)
}

func (s *Step) update() error {
	// alternate between SPLIT and MERGE
	s.Type = s.Type ^ 1

	// new From is where the coins were last sent to
	s.From = append([]string{}, s.To...)

	// get amount to divide across s.To (sum of s.From balances)
	amount, err := s.getAmount(s.From)
	if err != nil {
		return err
	}
	s.Amount = amount

	// reset for adding random addresses below
	s.To = make([]string, 0)

searching:
	for {
		// get a random index of Addrs
		i := rand.Intn(len(s.Addrs))

		// skip if the address is already in From
		for _, e := range s.From {
			if s.Addrs[i] == e {
				continue searching
			}
		}

		// skip if the address is already in To
		for _, e := range s.To {
			if s.Addrs[i] == e {
				continue searching
			}
		}

		s.To = append(s.To, s.Addrs[i])

		if (len(s.To) == 1 && s.Type == MERGE) ||
			(len(s.To) == SIZE && s.Type == SPLIT) {
			return nil
		}
	}
}

func (s *Step) Run() (*StepResult, error) {
	err := s.update()
	if err != nil {
		return nil, err
	}

	// convert []string to []cli.SendAmount
	to := make([]cli.SendAmount, len(s.To), len(s.To))
	for i := range s.To {
		to[i] = cli.SendAmount{s.To[i], s.Amount / uint64(len(s.To))}
	}

	start := time.Now()

	println("creating...")

	// create transaction
	tx, err := cli.CreateRawTx(s.Client, s.Wallet, s.From, s.From[0], to)
	if err != nil {
		return nil, err
	}

	println("injecting...")

	// inject transaction
	txId, err := s.Client.InjectTransaction(tx)
	if err != nil {
		return nil, err
	}

	println("waiting...")

	// track transaction status
	status, err := s.Wait(txId)
	if err != nil {
		return nil, err
	}

	end := time.Now()

	return &StepResult{
		Id:       txId,
		Type:     s.Type,
		From:     s.From,
		To:       to,
		Start:    start,
		End:      end,
		Duration: end.Sub(start),
		Status:   status,
	}, nil
}

func (s *Step) Cleanup() (*StepResult, error) {
	// send from all adresses not including the destination
	from := append([]string{}, s.Addrs[1:]...)

	// get the sum of all from addresses
	amount, err := s.getAmount(from)
	if err != nil {
		return nil, err
	}

	// send everything to the first address generated by seed
	to := []cli.SendAmount{{s.Addrs[0], amount}}

	start := time.Now()

	// create tx
	tx, err := cli.CreateRawTx(s.Client, s.Wallet, from, s.Addrs[0], to)
	if err != nil {
		return nil, err
	}

	// inject tx
	txId, err := s.Client.InjectTransaction(tx)
	if err != nil {
		return nil, err
	}

	// wait for confirmation
	status, err := s.Wait(txId)
	if err != nil {
		return nil, err
	}

	end := time.Now()

	return &StepResult{
		Id:       txId,
		Type:     CLEAN,
		From:     from,
		To:       to,
		Start:    start,
		End:      end,
		Duration: end.Sub(start),
		Status:   status,
	}, nil
}

func (s *Step) Wait(txId string) (*visor.TransactionStatus, error) {
	var (
		tx  *webrpc.TxnResult
		err error
	)

	for {
		if tx, err = s.Client.GetTransactionByID(txId); err != nil {
			println(err)
			continue
		}

		if tx.Transaction.Status.Confirmed {
			return &tx.Transaction.Status, nil
		}

		// once every 100 milliseconds
		<-time.After(time.Second / 10)
	}
}
