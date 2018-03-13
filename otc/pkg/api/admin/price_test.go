package admin

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
)

func MockPriceSend(data string) (*currencies.Currencies, string) {
	curs := MockCurrencies()

	res := httptest.NewRecorder()
	req := MockRequest(data)

	Price(curs, nil)(res, req)

	out, _ := ioutil.ReadAll(res.Body)
	return curs, strings.TrimSpace(string(out))
}

func TestPriceBad(t *testing.T) {
	_, res := MockPriceSend("bad json")

	if res != "invalid JSON" {
		t.Fatalf(`expected "invalid JSON", got "%s"`, res)
	}
}

func TestPriceGood(t *testing.T) {
	curs, res := MockPriceSend(`{"price":100000000}`)

	if res != "" {
		t.Fatalf(`expected empty response, got "%s"`, res)
	}

	if curs.Prices[otc.BTC].Sources[currencies.INTERNAL].Amount != 1e8 {
		t.Fatal("price wasn't set")
	}
}
