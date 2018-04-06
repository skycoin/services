package sender

import (
	"fmt"
	"testing"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
)

type Mock struct {
	Fail bool
}

func (m *Mock) Balance(string) (uint64, error) { return 0, nil }
func (m *Mock) Confirmed(string) (bool, error) { return false, nil }
func (m *Mock) Address() (string, error)       { return "", nil }
func (m *Mock) Used() ([]string, error)        { return nil, nil }
func (m *Mock) Connected() (bool, error)       { return false, nil }
func (m *Mock) Holding() (uint64, error)       { return 0, nil }
func (m *Mock) Stop() error                    { return nil }

func (m *Mock) Send(string, uint64) (string, error) {
	if m.Fail {
		return "", fmt.Errorf("fail!")
	}
	return "txid", nil
}

func TestTask(t *testing.T) {
	curs := &currencies.Currencies{
		Prices: map[otc.Currency]*currencies.Pricer{
			otc.BTC: &currencies.Pricer{
				Using: currencies.INTERNAL,
				Sources: map[currencies.Source]*currencies.Price{
					currencies.INTERNAL: currencies.NewPrice(200000),
				},
			},
		},
		Connections: map[otc.Currency]currencies.Connection{
			otc.SKY: &Mock{false},
		},
	}

	work := &otc.Work{
		Order: &otc.Order{
			User: &otc.User{
				Drop: &otc.Drop{
					Address:  "address",
					Currency: otc.BTC,
				},
			},
			Amount: 100000000,
			Times:  &otc.Times{},
		},
		Done: make(chan *otc.Result, 1),
	}

	if _, err := Task(curs)(work); err != nil {
		t.Fatal(err)
	}

	if work.Order.Purchase.TxId != "txid" {
		t.Fatal("didn't complete purchase")
	}

	if work.Order.Status != otc.CONFIRM {
		t.Fatal("didn't change order status")
	}
}

func TestTaskBadPrice(t *testing.T) {
	curs := &currencies.Currencies{
		Connections: map[otc.Currency]currencies.Connection{
			otc.SKY: &Mock{false},
		},
	}

	work := &otc.Work{
		Order: &otc.Order{
			User: &otc.User{
				Drop: &otc.Drop{
					Address:  "address",
					Currency: otc.BTC,
				},
			},
			Amount: 100000000,
			Times:  &otc.Times{},
		},
		Done: make(chan *otc.Result, 1),
	}

	if _, err := Task(curs)(work); err == nil {
		t.Fatal("should've returned an error")
	}
}

func TestTaskBadConnection(t *testing.T) {
	curs := &currencies.Currencies{
		Prices: map[otc.Currency]*currencies.Pricer{
			otc.BTC: &currencies.Pricer{
				Using: currencies.INTERNAL,
				Sources: map[currencies.Source]*currencies.Price{
					currencies.INTERNAL: currencies.NewPrice(200000),
				},
			},
		},
		Connections: map[otc.Currency]currencies.Connection{
			otc.SKY: &Mock{true},
		},
	}

	work := &otc.Work{
		Order: &otc.Order{
			User: &otc.User{
				Drop: &otc.Drop{
					Address:  "address",
					Currency: otc.BTC,
				},
			},
			Amount: 100000000,
			Times:  &otc.Times{},
		},
		Done: make(chan *otc.Result, 1),
	}

	if _, err := Task(curs)(work); err == nil {
		t.Fatal("should've returned an error")
	}
}
