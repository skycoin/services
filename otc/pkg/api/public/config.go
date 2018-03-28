package public

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/model"
	"github.com/skycoin/services/otc/pkg/otc"
)

func Config(curs *currencies.Currencies, modl *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		holding, err := curs.Holding(otc.SKY)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		price, err := curs.Price(otc.BTC)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		var status string

		if modl.Controller.Paused() {
			status = "PAUSED"
		} else {
			status = "WORKING"
		}

		json.NewEncoder(w).Encode(&struct {
			Status string `json:"otcStatus"`
			// TODO: change to holding
			Holding uint64 `json:"balance"`
			Price   uint64 `json:"price"`
		}{status, holding, price})
	}
}
