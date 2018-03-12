package currencies

import (
	"errors"
	"time"

	"github.com/skycoin/services/otc/pkg/exchange"
	"github.com/skycoin/services/otc/pkg/otc"
)

var (
	ErrConnExists   error = errors.New("connection already exists")
	ErrConnMissing  error = errors.New("connection missing")
	ErrPriceMissing error = errors.New("price missing")
	ErrZeroAmount   error = errors.New("zero amount")
)

type Connection interface {
	Balance(string) (uint64, error)
	Confirmed(string) (bool, error)
	Send(string, uint64) (string, error)
	Address() (string, error)
	Connected() (bool, error)
	Holding() (uint64, error)
	Stop() error
}

type Currencies struct {
	Prices      map[otc.Currency]*Pricer
	Connections map[otc.Currency]Connection
}

func New() *Currencies {
	return &Currencies{
		Prices:      make(map[otc.Currency]*Pricer),
		Connections: make(map[otc.Currency]Connection),
	}
}

func (c *Currencies) Add(curr otc.Currency, conn Connection) error {
	if c.Connections[curr] != nil {
		return ErrConnExists
	}

	c.Connections[curr] = conn

	if curr == otc.BTC {
		c.Prices[curr] = &Pricer{
			Using: INTERNAL,
			Sources: map[Source]*Price{
				INTERNAL: NewPrice(200000),
			},
		}

		go func() {
			for {
				price, err := exchange.GetBTCValue()
				if err != nil {
					c.Prices[curr].SetSource(INTERNAL)
				} else {
					c.Prices[curr].SetPrice(EXCHANGE, price)
					c.Prices[curr].SetSource(EXCHANGE)
				}

				<-time.After(time.Minute)
			}
		}()
	}

	return nil
}

func (c *Currencies) Holding(curr otc.Currency) (uint64, error) {
	if c.Connections[curr] == nil {
		return 0, ErrConnMissing
	}

	return c.Connections[curr].Holding()
}

func (c *Currencies) Balance(drop *otc.Drop) (uint64, error) {
	if c.Connections[drop.Currency] == nil {
		return 0, ErrConnMissing
	}

	return c.Connections[drop.Currency].Balance(drop.Address)
}

func (c *Currencies) Value(drop *otc.Drop) (uint64, string, error) {
	if c.Prices[drop.Currency] == nil {
		return 0, "", ErrPriceMissing
	}

	if drop.Amount == 0 {
		return 0, "", ErrZeroAmount
	}

	price, source, _ := c.Prices[drop.Currency].GetPrice()
	return uint64(float64(float64(drop.Amount)/float64(price)*1e2)) * 1e4, string(source), nil
}

func (c *Currencies) Send(curr otc.Currency, addr string, amount uint64) (string, error) {
	if c.Connections[curr] == nil {
		return "", ErrConnMissing
	}

	if amount == 0 {
		return "", ErrZeroAmount
	}

	return c.Connections[curr].Send(addr, amount)
}

func (c *Currencies) Confirmed(curr otc.Currency, txid string) (bool, error) {
	if c.Connections[curr] == nil {
		return false, ErrConnMissing
	}

	return c.Connections[curr].Confirmed(txid)
}

func (c *Currencies) Address(curr otc.Currency) (string, error) {
	if c.Connections[curr] == nil {
		return "", ErrConnMissing
	}

	return c.Connections[curr].Address()
}

func (c *Currencies) Price(curr otc.Currency) (uint64, error) {
	if c.Prices[curr] == nil {
		return 0, ErrPriceMissing
	}

	price, _, _ := c.Prices[curr].GetPrice()
	return price, nil
}

func (c *Currencies) Source(curr otc.Currency) (Source, error) {
	if c.Prices[curr] == nil {
		return "", ErrPriceMissing
	}

	source := c.Prices[curr].GetSource()
	return source, nil
}
