package admin

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/model"
	"github.com/skycoin/services/otc/pkg/otc"
)

func Price(curs *currencies.Currencies, modl *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			req = &struct {
				Price uint64 `json:"price"`
			}{}
			err error
		)

		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		curs.Prices[otc.BTC].Sources[currencies.INTERNAL].Set(req.Price)
	}
}
