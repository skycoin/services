package scanner

import (
	"fmt"
	"testing"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
)

type MockConnection struct{}

func (c *MockConnection) Balance(addr string) (uint64, error) {
	if addr == "empty" {
		return 0, nil
	}

	if addr == "full" {
		return 1000, nil
	}

	if addr == "error" {
		return 0, fmt.Errorf("")
	}

	return 0, nil
}

func (c *MockConnection) Confirmed(txid string) (bool, error) {
	return false, nil
}

func (c *MockConnection) Send(addr string, amount uint64) (string, error) {
	return "txid", nil
}

func (c *MockConnection) Address() (string, error) {
	return "mock", nil
}

func (c *MockConnection) Connected() (bool, error) {
	return false, nil
}

func (c *MockConnection) Holding() (uint64, error) {
	return 0, nil
}

func (c *MockConnection) Stop() error {
	return nil
}

func TestTask(t *testing.T) {
	curs := currencies.New()
	curs.Add(otc.BTC, &MockConnection{})

	var (
		remove bool
		err    error
	)

	if remove, err = Task(curs)(&otc.Work{
		Request: &otc.Request{
			Times: &otc.Times{},
			Drop: &otc.Drop{
				Currency: otc.BTC,
				Address:  "empty",
			},
		},
	}); remove || err != nil {
		t.Fatal("shouldn't be removed")
	}

	if remove, err = Task(curs)(&otc.Work{
		Request: &otc.Request{
			Times: &otc.Times{},
			Drop: &otc.Drop{
				Currency: otc.BTC,
				Address:  "full",
			},
		},
	}); !remove || err != nil {
		t.Fatal("should be removed")
	}

	if remove, err = Task(curs)(&otc.Work{
		Request: &otc.Request{
			Times: &otc.Times{},
			Drop: &otc.Drop{
				Currency: otc.BTC,
				Address:  "error",
			},
		},
	}); err == nil {
		t.Fatal("should be error")
	}

	work := &otc.Work{
		Request: &otc.Request{
			Times: &otc.Times{},
			Drop: &otc.Drop{
				Currency: otc.BTC,
				Address:  "full",
			},
		},
	}

	if remove, err = Task(curs)(work); err != nil {
		t.Fatal(err)
	}

	if work.Request.Status != otc.SEND {
		t.Fatal("status should be SEND")
	}

	if work.Request.Drop.Amount != 1000 {
		t.Fatal("amount should be set")
	}

	if work.Request.Times.DepositedAt == 0 {
		t.Fatal("deposited at time should be set")
	}
}
