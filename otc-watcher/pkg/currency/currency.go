package currency

import (
	"errors"

	"github.com/skycoin/services/otc/pkg/otc"
)

var (
	ErrConnMissing = errors.New("connection missing")
)

type Connection interface {
	Stop() error
	Scan(uint64) (chan *otc.Block, error)
	Get(uint64) (*otc.Block, error)
	Height() (uint64, error)
}

type Connections map[otc.Currency]Connection

func (c Connections) Get(cur otc.Currency, height uint64) (*otc.Block, error) {
	if c[cur] == nil {
		return nil, ErrConnMissing
	}

	return c[cur].Get(height)
}

func (c Connections) Height(cur otc.Currency) (uint64, error) {
	if c[cur] == nil {
		return 0, ErrConnMissing
	}

	return c[cur].Height()
}

func (c Connections) Scan(cur otc.Currency, from uint64) (chan *otc.Block, error) {
	if c[cur] == nil {
		return nil, ErrConnMissing
	}

	return c[cur].Scan(from)
}
