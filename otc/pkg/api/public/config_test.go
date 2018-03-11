package public

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
)

func TestConfig(t *testing.T) {
	curs_one := &currencies.Currencies{
		Connections: map[otc.Currency]currencies.Connection{
			otc.SKY: &MockConnection{Bad: true},
		},
	}

	curs_two := &currencies.Currencies{
		Connections: map[otc.Currency]currencies.Connection{
			otc.SKY: &MockConnection{},
		},
		Prices: map[otc.Currency]*currencies.Pricer{
			otc.SKY: nil,
		},
	}

	curs_three := &currencies.Currencies{
		Connections: map[otc.Currency]currencies.Connection{
			otc.SKY: &MockConnection{},
		},
		Prices: map[otc.Currency]*currencies.Pricer{
			otc.BTC: &currencies.Pricer{
				Using: currencies.INTERNAL,
				Sources: map[currencies.Source]*currencies.Price{
					currencies.INTERNAL: currencies.NewPrice(100),
				},
			},
		},
	}

	tests := map[*currencies.Currencies]string{
		curs_one:   "server error",
		curs_two:   "server error",
		curs_three: `{"otcStatus":"WORKING","balance":0,"price":100}`,
	}

	for curs, expected := range tests {
		req := httptest.NewRequest("GET", "http:///", nil)
		res := httptest.NewRecorder()

		Config(curs, nil)(res, req)

		out, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		if strings.TrimSpace(string(out)) != expected {
			t.Fatalf(`expected "%s", got "%s"`, expected,
				strings.TrimSpace(string(out)))
		}
	}
}
