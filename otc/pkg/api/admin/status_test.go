package admin

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/model"
)

func MockStatusSend(curs *currencies.Currencies) (*model.Model, string) {
	modl := MockModel()

	res := httptest.NewRecorder()
	req := MockRequest("")

	Status(curs, modl)(res, req)

	out, _ := ioutil.ReadAll(res.Body)
	return modl, strings.TrimSpace(string(out))
}

func TestStatusBad(t *testing.T) {
	_, res := MockStatusSend(&currencies.Currencies{})

	if res != "server error" {
		t.Fatalf(`expected "server error", got "%s"`, res)
	}
}

func TestStatusGood(t *testing.T) {
	_, res := MockStatusSend(MockCurrencies())

	if res[0] != '{' {
		t.Fatalf(`expected JSON output, got "%s"`, res)
	}
}
