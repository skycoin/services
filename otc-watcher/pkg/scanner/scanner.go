package scanner

import (
	"errors"

	"log"

	"github.com/skycoin/services/otc/pkg/otc"

	"github.com/skycoin/services/otc-watcher/pkg/currency"
)

var (
	ErrAddressMissing = errors.New("address not in watch list")
)

type Scanner struct {
	Scanning map[otc.Currency]*Storage
}

func New(cons currency.Connections) (*Scanner, error) {
	s := &Scanner{make(map[otc.Currency]*Storage, 0)}

	// load from disk or create
	if err := s.Load(cons); err != nil {
		return nil, err
	}

	// for each connection (supported currency)
	for cur, con := range cons {
		// get blocks chan from connection
		blocks, err := con.Scan(s.Scanning[cur].Updated.Height)
		if err != nil {
			return nil, err
		}
		// start scanning
		go s.Scan(cur, blocks)
	}

	return s, nil
}

// TODO: stop connections and write to disk
func (s *Scanner) Stop() error {
	return nil
}

func (s *Scanner) Scan(cur otc.Currency, blocks chan *otc.Block) {
	for {
		// get block from connection channel
		block := <-blocks
		// TODO: use logger
		log.Printf("scanning block %d\n", block.Height)

		// update storage based on received block
		s.Scanning[cur].Update(block)

		// TODO: handle error better
		//
		// TODO: move Save function to *Storage so there's no need to pass
		//       cur as a param
		if err := s.Save(cur); err != nil {
			println(err.Error())
		}
	}
}

func (s *Scanner) Register(drop *otc.Drop) error {
	// check that connection exists
	if s.Scanning[drop.Currency] == nil {
		return currency.ErrConnMissing
	}

	// add to storage
	s.Scanning[drop.Currency].Register(drop.Address)

	// save to disk
	return s.Save(drop.Currency)
}

func (s *Scanner) Outputs(drop *otc.Drop) (otc.Outputs, error) {
	// check that connection exists
	if s.Scanning[drop.Currency] == nil {
		return nil, currency.ErrConnMissing
	}

	// check that address is registered
	if s.Scanning[drop.Currency].Addresses[drop.Address] == nil {
		return nil, ErrAddressMissing
	}

	// get outputs from storage
	return s.Scanning[drop.Currency].Outputs(drop.Address), nil
}
