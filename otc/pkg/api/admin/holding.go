package admin

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/model"
	"github.com/skycoin/services/otc/pkg/otc"
)

func Holding(curr otc.Currency, curs *currencies.Currencies, modl *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		holding, err := curs.Holding(curr)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(&struct {
			Holding uint64 `json:"holding"`
		}{holding})
	}
}
