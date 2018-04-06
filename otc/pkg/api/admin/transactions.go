package admin

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/model"
	"github.com/skycoin/services/otc/pkg/otc"
)

func Transactions(curs *currencies.Currencies, modl *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all := modl.Orders()

		sort.Slice(all, func(i, j int) bool {
			return all[i].Times.CreatedAt > all[j].Times.CreatedAt
		})

		json.NewEncoder(w).Encode(&all)
	}
}

func TransactionsPending(curs *currencies.Currencies, modl *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all := modl.Orders()
		pending := make([]otc.Order, 0)

		for _, order := range all {
			if order.Status != otc.DONE {
				pending = append(pending, order)
			}
		}

		sort.Slice(pending, func(i, j int) bool {
			return pending[i].Times.CreatedAt > pending[j].Times.CreatedAt
		})

		json.NewEncoder(w).Encode(&pending)
	}
}

func TransactionsCompleted(curs *currencies.Currencies, modl *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all := modl.Orders()
		completed := make([]otc.Order, 0)

		for _, order := range all {
			if order.Status == otc.DONE {
				completed = append(completed, order)
			}
		}

		sort.Slice(completed, func(i, j int) bool {
			return completed[i].Times.CreatedAt > completed[j].Times.CreatedAt
		})

		json.NewEncoder(w).Encode(&completed)
	}
}
