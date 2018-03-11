package monitor

import (
	"fmt"
	"testing"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
)

type MockConnection struct{}

func (c *MockConnection) Balance(addr string) (uint64, error) {
	return 0, nil
}

func (c *MockConnection) Confirmed(txid string) (bool, error) {
	if txid == "confirmed" {
		return true, nil
	}

	if txid == "unconfirmed" {
		return false, nil
	}

	if txid == "error" {
		return false, fmt.Errorf("error")
	}

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
	curs.Add(otc.SKY, &MockConnection{})

	var (
		remove bool
		err    error
	)

	if remove, err = Task(curs)(&otc.Work{
		Request: &otc.Request{TxId: "confirmed", Times: &otc.Times{}},
	}); !remove || err != nil {
		t.Fatal("should be removed")
	}

	if remove, err = Task(curs)(&otc.Work{
		Request: &otc.Request{TxId: "unconfirmed", Times: &otc.Times{}},
	}); remove || err != nil {
		t.Fatal("shouldn't be removed")
	}

	if remove, err = Task(curs)(&otc.Work{
		Request: &otc.Request{TxId: "error", Times: &otc.Times{}},
	}); err == nil {
		t.Fatal("should be an error")
	}

	work := &otc.Work{
		Request: &otc.Request{
			TxId:  "confirmed",
			Times: &otc.Times{},
		},
	}

	if remove, err = Task(curs)(work); err != nil {
		t.Fatal(err)
	}

	if work.Request.Status != otc.DONE {
		t.Fatal("monitor should set status to DONE")
	}

	if work.Request.Times.ConfirmedAt == 0 {
		t.Fatal("monitor didn't set ConfirmedAt time")
	}
}
