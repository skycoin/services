package public

import (
	"net/http/httptest"
	"testing"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/otc"
)

func TestAdminNew(t *testing.T) {
	mux := New(&currencies.Currencies{
		Connections: map[otc.Currency]currencies.Connection{
			otc.SKY: &MockConnection{},
		},
		Prices: map[otc.Currency]*currencies.Pricer{
			otc.SKY: &currencies.Pricer{
				Using: currencies.INTERNAL,
				Sources: map[currencies.Source]*currencies.Price{
					currencies.INTERNAL: currencies.NewPrice(100),
				},
			},
		},
	}, nil)

	paths := []string{
		"/api/bind",
		"/api/status",
		"/api/config",
	}

	for _, path := range paths {
		req := httptest.NewRequest("GET", path, nil)
		res := httptest.NewRecorder()
		mux.ServeHTTP(res, req)

		if res.Result().StatusCode == 404 {
			t.Fatalf("%s returning 404", path)
		}
	}
}
