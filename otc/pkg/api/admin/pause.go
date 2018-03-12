package admin

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/model"
)

func Pause(curs *currencies.Currencies, modl *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			req = &struct {
				Pause bool `json:"pause"`
			}{}
			err error
		)

		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if req.Pause {
			modl.Pause()
		} else {
			modl.Unpause()
		}
	}
}
