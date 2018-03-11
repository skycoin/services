package admin

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/model"
	"github.com/skycoin/services/otc/pkg/otc"
)

func Source(curs *currencies.Currencies, modl *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			req = &struct {
				Source string `json:"source"`
			}{}
			err error
		)

		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		var source currencies.Source

		// check that source option exists
		if req.Source != "exchange" && req.Source != "internal" {
			http.Error(w, "invalid price source", http.StatusBadRequest)
			return
		} else if req.Source == "exchange" {
			source = currencies.EXCHANGE
		} else if req.Source == "internal" {
			source = currencies.INTERNAL
		}

		curs.Prices[otc.BTC].SetSource(source)
	}
}
