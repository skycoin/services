package admin

import (
	"net/http/httptest"
	"testing"
)

func TestAdminNew(t *testing.T) {
	mux := New(MockCurrencies(), MockModel())
	paths := []string{
		"/api/status",
		"/api/pause",
		"/api/price",
		"/api/source",
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
