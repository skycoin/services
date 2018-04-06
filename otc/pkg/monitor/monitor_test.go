package monitor

import (
	"fmt"
	"testing"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
)

type Mock struct {
	Type string
}

func (m *Mock) Balance(string) (uint64, error)      { return 0, nil }
func (m *Mock) Address() (string, error)            { return "", nil }
func (m *Mock) Used() ([]string, error)             { return nil, nil }
func (m *Mock) Connected() (bool, error)            { return false, nil }
func (m *Mock) Holding() (uint64, error)            { return 0, nil }
func (m *Mock) Stop() error                         { return nil }
func (m *Mock) Send(string, uint64) (string, error) { return "", nil }

func (m *Mock) Confirmed(string) (bool, error) {
	if m.Type == "confirmed" {
		return true, nil
	} else if m.Type == "unconfirmed" {
		return false, nil
	} else if m.Type == "error" {
		return false, fmt.Errorf("error!")
	}
	return false, nil
}

func TestTask(t *testing.T) {
	curs := &currencies.Currencies{
		Connections: map[otc.Currency]currencies.Connection{
			otc.SKY: &Mock{"confirmed"},
		},
	}

	work := &otc.Work{
		Order: &otc.Order{
			Purchase: &otc.Purchase{
				TxId: "txid",
			},
			Times: &otc.Times{},
		},
		Done: make(chan *otc.Result, 1),
	}

	done, err := Task(curs)(work)

	if err != nil {
		t.Fatal("shouldn't be an error")
	}

	if !done {
		t.Fatal("should be done")
	}
}

func TestTaskUnconfirmed(t *testing.T) {
	curs := &currencies.Currencies{
		Connections: map[otc.Currency]currencies.Connection{
			otc.SKY: &Mock{"unconfirmed"},
		},
	}

	work := &otc.Work{
		Order: &otc.Order{
			Purchase: &otc.Purchase{
				TxId: "txid",
			},
			Times: &otc.Times{},
		},
		Done: make(chan *otc.Result, 1),
	}

	done, err := Task(curs)(work)

	if err != nil {
		t.Fatal("shouldn't be an error")
	}

	if done {
		t.Fatal("shouldn't be done")
	}
}

func TestTaskError(t *testing.T) {
	curs := &currencies.Currencies{
		Connections: map[otc.Currency]currencies.Connection{
			otc.SKY: &Mock{"error"},
		},
	}

	work := &otc.Work{
		Order: &otc.Order{
			Purchase: &otc.Purchase{
				TxId: "txid",
			},
			Times: &otc.Times{},
		},
		Done: make(chan *otc.Result, 1),
	}

	done, err := Task(curs)(work)

	if err == nil {
		t.Fatal("should be an error")
	}

	if !done {
		t.Fatal("should be done")
	}
}
