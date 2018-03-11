package sender

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
	return false, nil
}

func (c *MockConnection) Send(addr string, amount uint64) (string, error) {
	if addr == "good" {
		return "txid", nil
	}

	if addr == "bad" {
		return "", fmt.Errorf("")
	}

	return "", nil
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

	curs.Prices[otc.BTC] = &currencies.Pricer{
		Using: currencies.INTERNAL,
		Sources: map[currencies.Source]*currencies.Price{
			currencies.INTERNAL: currencies.NewPrice(1000),
		},
	}

	var (
		remove bool
		err    error
	)

	if remove, err = Task(curs)(&otc.Work{
		Request: &otc.Request{
			Drop: &otc.Drop{
				Address:  "",
				Currency: otc.ETH,
				Amount:   1000,
			},
		},
	}); !remove || err == nil {
		t.Fatal("should remove with error")
	}

	if remove, err = Task(curs)(&otc.Work{
		Request: &otc.Request{
			Address: "bad",
			Times:   &otc.Times{},
			Drop: &otc.Drop{
				Address:  "",
				Currency: otc.BTC,
				Amount:   1000,
			},
		},
	}); !remove || err == nil {
		t.Fatal("should remove with error")
	}

	work := &otc.Work{
		Request: &otc.Request{
			Address: "good",
			Times:   &otc.Times{},
			Drop: &otc.Drop{
				Address:  "",
				Currency: otc.BTC,
				Amount:   1000,
			},
		},
	}

	if remove, err = Task(curs)(work); err != nil {
		t.Fatal(err)
	}

	if work.Request.Status != otc.CONFIRM {
		t.Fatal("status should be CONFIRM")
	}

	if work.Request.TxId != "txid" {
		t.Fatal("txid not set")
	}

	if work.Request.Times.SentAt == 0 {
		t.Fatal("sent at time not set")
	}
}
