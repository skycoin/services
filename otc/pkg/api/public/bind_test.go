package public

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/skycoin/services/otc/pkg/actor"
	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/model"
	"github.com/skycoin/services/otc/pkg/otc"
)

type MockConnection struct {
	Bad bool
}

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
	return "", nil
}

func (c *MockConnection) Address() (string, error) {
	if c.Bad {
		return "", fmt.Errorf("bad")
	}

	return "mock", nil
}

func (c *MockConnection) Connected() (bool, error) {
	return false, nil
}

func (c *MockConnection) Holding() (uint64, error) {
	if c.Bad {
		return 0, fmt.Errorf("bad")
	}

	return 0, nil
}

func (c *MockConnection) Stop() error {
	return nil
}

func TestBind(t *testing.T) {
	curs := &currencies.Currencies{
		Prices: map[otc.Currency]*currencies.Pricer{
			otc.BTC: &currencies.Pricer{
				Using: currencies.INTERNAL,
				Sources: map[currencies.Source]*currencies.Price{
					currencies.INTERNAL: currencies.NewPrice(100),
					currencies.EXCHANGE: currencies.NewPrice(200),
				},
			},
		},
		Connections: map[otc.Currency]currencies.Connection{
			otc.BTC: &MockConnection{},
			otc.SKY: &MockConnection{},
			otc.ETH: &MockConnection{Bad: true},
		},
	}
	modl := &model.Model{
		Running: true,
		Lookup:  make(map[string]*otc.Request),
		Logger:  log.New(ioutil.Discard, "", 0),
		Router:  actor.New(nil, nil),
	}

	tests := [][]string{
		{
			`bad json`,
			`invalid JSON`,
		},
		{
			`{"address":"234234",
			  "drop_currency":"BTC"}`,
			`invalid skycoin address`,
		},
		{
			`{"address":"2dvVgeKNU7UHdvvBUVZXbBaxoTkpemo1cmg",
			  "drop_currency":"???"}`,
			`not supported`,
		},
		{
			`{"address":"2dvVgeKNU7UHdvvBUVZXbBaxoTkpemo1cmg",
			  "drop_currency":"ETH"}`,
			`server error`,
		},
		{
			`{"address":"2dvVgeKNU7UHdvvBUVZXbBaxoTkpemo1cmg",
			  "drop_currency":"SKY"}`,
			`server error`,
		},
		{
			`{"address":"2dvVgeKNU7UHdvvBUVZXbBaxoTkpemo1cmg",
			  "drop_currency":"BTC"}`,
			`{"drop_address":"mock","drop_currency":"BTC","drop_value":100}`,
		},
	}

	for _, test := range tests {
		var buf bytes.Buffer
		buf.WriteString(test[0])
		req := httptest.NewRequest("GET", "http:///", &buf)
		res := httptest.NewRecorder()

		Bind(curs, modl)(res, req)

		out, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		if strings.TrimSpace(string(out)) != test[1] {
			t.Fatalf(`expected "%s", got "%s"`, test[1],
				strings.TrimSpace(string(out)))
		}
	}
}
