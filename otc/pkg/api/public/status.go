package public

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/model"
)

func Status(curs *currencies.Currencies, modl *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			data struct {
				DropAddress  string `json:"drop_address"`
				DropCurrency string `json:"drop_currency"`
			}
			err error
		)

		if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		user, err := modl.Lookup.GetStatus(
			data.DropCurrency + ":" + data.DropAddress,
		)
		if err != nil {
			http.Error(w, "user missing", http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(user.Orders)
	}
}
