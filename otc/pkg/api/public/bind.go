package public

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/model"
	"github.com/skycoin/services/otc/pkg/otc"
	"github.com/skycoin/skycoin/src/cipher"
)

func Bind(curs *currencies.Currencies, modl *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			data struct {
				Affiliate    string `json:"affiliate"`
				Address      string `json:"address"`
				DropCurrency string `json:"drop_currency"`
			}
			err error
		)

		if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if modl.Paused() {
			http.Error(w, "paused", http.StatusInternalServerError)
			return
		}

		curr := otc.Currency(data.DropCurrency)

		addr, err := cipher.DecodeBase58Address(data.Address)
		if err != nil {
			http.Error(w, "invalid skycoin address", http.StatusBadRequest)
			return
		}

		dropAddr, err := curs.Address(curr)
		if err != nil {
			if err == currencies.ErrConnMissing {
				http.Error(w, "not supported", http.StatusBadRequest)
			} else {
				println(err.Error())
				http.Error(w, "server error", http.StatusInternalServerError)
			}
			return
		}

		req := &otc.Request{
			Affiliate: data.Affiliate,
			Address:   addr.String(),
			Status:    otc.NEW,
			TxId:      "",
			Times: &otc.Times{
				CreatedAt: time.Now().UTC().Unix(),
			},
			Drop: &otc.Drop{
				Address:  dropAddr,
				Currency: curr,
				Amount:   0,
			},
		}

		price, err := curs.Price(curr)
		if err != nil {
			println(err.Error())
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		modl.Add(req)

		json.NewEncoder(w).Encode(&struct {
			DropAddress  string       `json:"drop_address"`
			DropCurrency otc.Currency `json:"drop_currency"`
			// TODO: change to price
			DropValue uint64 `json:"drop_value"`
		}{dropAddr, curr, price})
	}
}
