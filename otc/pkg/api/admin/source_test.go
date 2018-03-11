package admin

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
)

func MockCurrencies() *currencies.Currencies {
	return &currencies.Currencies{
		Prices: map[otc.Currency]*currencies.Pricer{
			otc.BTC: &currencies.Pricer{
				Using: currencies.INTERNAL,
				Sources: map[currencies.Source]*currencies.Price{
					currencies.INTERNAL: currencies.NewPrice(100),
					currencies.EXCHANGE: currencies.NewPrice(200),
				},
			},
		},
	}
}

func MockRequest(data string) *http.Request {
	var body bytes.Buffer
	body.WriteString(data)
	return httptest.NewRequest("GET", "http:///", &body)
}

func MockSourceSend(data string) (*currencies.Currencies, string) {
	curs := MockCurrencies()

	res := httptest.NewRecorder()
	req := MockRequest(data)

	Source(curs, nil)(res, req)

	out, _ := ioutil.ReadAll(res.Body)
	return curs, strings.TrimSpace(string(out))
}

func TestSourceInvalidJSON(t *testing.T) {
	_, res := MockSourceSend("bad json")

	if res != "invalid JSON" {
		t.Fatalf(`expected "invalid JSON", got "%s"`, res)
	}
}

func TestSourceInvalidSource(t *testing.T) {
	_, res := MockSourceSend(`{"source":"bad source"}`)

	if res != "invalid price source" {
		t.Fatalf(`expected "invalid price source", got "%s"`, res)
	}
}

func TestSourceExchange(t *testing.T) {
	curs, res := MockSourceSend(`{"source":"exchange"}`)

	if res != "" {
		t.Fatalf(`expected empty response, got "%s"`, res)
	}

	if curs.Prices[otc.BTC].Using != currencies.EXCHANGE {
		t.Fatal("source didn't change")
	}
}

func TestSourceInternal(t *testing.T) {
	curs, res := MockSourceSend(`{"source":"internal"}`)

	if res != "" {
		t.Fatalf(`expected empty response, got "%s"`, res)
	}

	if curs.Prices[otc.BTC].Using != currencies.INTERNAL {
		t.Fatal("source didn't change")
	}
}
