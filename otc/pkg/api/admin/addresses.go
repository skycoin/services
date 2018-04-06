package admin

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/model"
	"github.com/skycoin/services/otc/pkg/otc"
)

func Addresses(curr otc.Currency, curs *currencies.Currencies, modl *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type object struct {
			Address string `json:"address"`
			Balance uint64 `json:"balance"`
		}

		used, err := curs.Used(curr)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		addrs := make([]*object, len(used), len(used))

		for i := range used {
			balance, err := curs.Balance(&otc.Drop{used[i], curr})
			if err != nil {
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}

			addrs[i] = &object{
				Address: used[i],
				Balance: balance,
			}
		}

		json.NewEncoder(w).Encode(&addrs)
	}
}
