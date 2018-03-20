package currencies

import (
	"testing"

	"github.com/skycoin/services/otc/pkg/otc"
)

type MockConnection struct{}

func (c *MockConnection) Used() ([]string, error) {
	return nil, nil
}

func (c *MockConnection) Balance(addr string) (uint64, error) {
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

func TestCurrenciesNew(t *testing.T) {
	curs := New()

	if curs.Prices == nil {
		t.Fatal("nil prices")
	}

	if curs.Connections == nil {
		t.Fatal("nil connections")
	}
}

func TestCurrenciesAdd(t *testing.T) {
	curs := New()

	err := curs.Add(otc.BTC, &MockConnection{})
	if err != nil {
		t.Fatal(err)
	}

	err = curs.Add(otc.BTC, &MockConnection{})
	if err != ErrConnExists {
		t.Fatal(err)
	}

	if curs.Connections[otc.BTC] == nil {
		t.Fatal("add connection error")
	}

	if curs.Prices[otc.BTC] == nil {
		t.Fatal("add price error")
	}
}

func TestCurrenciesHolding(t *testing.T) {
	curs := New()
	curs.Add(otc.BTC, &MockConnection{})

	_, err := curs.Holding(otc.SKY)
	if err != ErrConnMissing {
		t.Fatal(err)
	}

	_, err = curs.Holding(otc.BTC)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCurrenciesBalance(t *testing.T) {
	curs := New()
	curs.Add(otc.BTC, &MockConnection{})

	_, err := curs.Balance(&otc.Drop{Currency: otc.SKY})
	if err != ErrConnMissing {
		t.Fatal(err)
	}

	_, err = curs.Balance(&otc.Drop{Currency: otc.BTC})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCurrenciesValue(t *testing.T) {
	curs := New()
	curs.Add(otc.BTC, &MockConnection{})
	curs.Prices[otc.BTC] = &Pricer{
		Using: INTERNAL,
		Sources: map[Source]*Price{
			INTERNAL: NewPrice(200000),
		},
	}

	_, _, err := curs.Value(&otc.Drop{
		Currency: otc.SKY,
	})
	if err != ErrPriceMissing {
		t.Fatal(err)
	}

	_, _, err = curs.Value(&otc.Drop{
		Currency: otc.BTC,
		Amount:   0,
	})
	if err != ErrZeroAmount {
		t.Fatal(err)
	}

	value, _, err := curs.Value(&otc.Drop{
		Currency: otc.BTC,
		Amount:   100000000,
	})
	if value != (500 * 1e6) {
		t.Fatal("bad value calculation")
	}
}

func TestCurrenciesSend(t *testing.T) {
	curs := New()
	curs.Add(otc.BTC, &MockConnection{})

	_, err := curs.Send(otc.SKY, "", 0)
	if err != ErrConnMissing {
		t.Fatal(err)
	}

	_, err = curs.Send(otc.BTC, "", 0)
	if err != ErrZeroAmount {
		t.Fatal(err)
	}

	txid, err := curs.Send(otc.BTC, "", 10)
	if err != nil {
		t.Fatal(err)
	}

	if txid != "txid" {
		t.Fatal("bad txid from send")
	}
}

func TestCurrenciesConfirmed(t *testing.T) {
	curs := New()
	curs.Add(otc.BTC, &MockConnection{})

	_, err := curs.Confirmed(otc.SKY, "")
	if err != ErrConnMissing {
		t.Fatal(err)
	}

	_, err = curs.Confirmed(otc.BTC, "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCurrenciesAddress(t *testing.T) {
	curs := New()
	curs.Add(otc.BTC, &MockConnection{})

	_, err := curs.Address(otc.SKY)
	if err != ErrConnMissing {
		t.Fatal(err)
	}

	_, err = curs.Address(otc.BTC)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCurrenciesPrice(t *testing.T) {
	curs := New()
	curs.Add(otc.BTC, &MockConnection{})
	curs.Prices[otc.BTC] = &Pricer{
		Using: INTERNAL,
		Sources: map[Source]*Price{
			INTERNAL: NewPrice(200000),
		},
	}

	_, err := curs.Price(otc.SKY)
	if err != ErrPriceMissing {
		t.Fatal(err)
	}

	_, err = curs.Price(otc.BTC)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCurrenciesSource(t *testing.T) {
	curs := New()
	curs.Add(otc.BTC, &MockConnection{})
	curs.Prices[otc.BTC] = &Pricer{
		Using: INTERNAL,
		Sources: map[Source]*Price{
			INTERNAL: NewPrice(200000),
		},
	}

	_, err := curs.Source(otc.SKY)
	if err != ErrPriceMissing {
		t.Fatal(err)
	}

	_, err = curs.Source(otc.BTC)
	if err != nil {
		t.Fatal(err)
	}
}
