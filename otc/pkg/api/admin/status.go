package admin

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/model"
	"github.com/skycoin/services/otc/pkg/otc"
)

func Status(curs *currencies.Currencies, modl *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			res = &struct {
				Prices struct {
					Internal        uint64 `json:"internal"`
					InternalUpdated int64  `json:"internal_updated"`
					Exchange        uint64 `json:"exchange"`
					ExchangeUpdated int64  `json:"exchange_updated"`
				} `json:"prices"`
				Source currencies.Source `json:"source"`
				Paused bool              `json:"paused"`
			}{}
			err error
		)

		res.Paused = modl.Controller.Paused()

		// TODO: add other currency support
		if res.Source, err = curs.Source(otc.BTC); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		// internal prices
		i, iu := curs.Prices[otc.BTC].Sources[currencies.INTERNAL].Get()
		res.Prices.Internal = i
		res.Prices.InternalUpdated = iu.UTC().Unix()

		// exchange prices
		e, eu := curs.Prices[otc.BTC].Sources[currencies.EXCHANGE].Get()
		res.Prices.Exchange = e
		res.Prices.ExchangeUpdated = eu.UTC().Unix()

		json.NewEncoder(w).Encode(&res)
	}
}
