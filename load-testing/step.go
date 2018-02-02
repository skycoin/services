package main

import (
	"errors"
	"log"
	"math/rand"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

type StepType int

const (
	// MERGE decribes a transaction from many addresses to 1.
	MERGE StepType = iota
	// SPLIT describes a transaction from 1 address to many.
	SPLIT
)

const (
	// SIZE is the number of addresses to SPLIT across. If SIZE is 2, then steps
	// will alternate between 1->2 and 2->1. If SIZE is 10, then steps will
	// alternate between 1->10 and 10->1.
	SIZE = 3
)

var (
	ErrNotEnoughAddrs = errors.New("not enough addresses to split across")
	ErrNoBalance      = errors.New("at least one address must have a spendable balance > 0")
)

type StepResult struct{}

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

func NewStep(c *webrpc.Client, w *wallet.Wallet, l *log.Logger, addrs []cipher.Address) (*Step, error) {
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

	// get initial s.To and s.Amount
	var err error
	s.To, s.Amount, err = s.findOrigin()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Step) findOrigin() ([]string, uint64, error) {
	unspent, err := s.Client.GetUnspentOutputs(s.Addrs)
	if err != nil {
		return nil, 0, err
	}

	spendable, err := unspent.Outputs.SpendableOutputs().ToUxArray()
	if err != nil {
		return nil, 0, err
	} else if len(spendable) == 0 {
		return nil, 0, ErrNoBalance
	}

	// just use the first address found
	origin := spendable[0].Body.Address.String()
	amount := spendable[0].Body.Coins / SIZE

	// find the id of initial address in the slice
	for _, addr := range s.Addrs {
		if addr == origin {
			return []string{addr}, amount, nil
		}
	}

	return nil, 0, ErrNoBalance
}

func (s *Step) update() {
	// alternate between SPLIT and MERGE
	s.Type = s.Type ^ 1

	// new From is where the coins were last sent to
	s.From = append([]string{}, s.To...)

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

		s.To = append(s.To, s.Addrs[i])

		if (len(s.To) == 1 && s.Type == MERGE) ||
			(len(s.To) == SIZE && s.Type == SPLIT) {
			return
		}
	}
}

func (s *Step) Run() *StepResult {
	s.update()

	if s.Type == MERGE {
		s.Logger.Println("MERGE", s.Amount)
	} else if s.Type == SPLIT {
		s.Logger.Println("SPLIT", s.Amount)
	}

	s.Logger.Printf("\tFROM\t%v\n", s.From)
	s.Logger.Printf("\tTO  \t%v\n", s.To)

	// create transaction
	// inject into network
	// track status

	return nil
}
