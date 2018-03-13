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
		all := modl.Reqs()

		sort.Slice(all, func(i, j int) bool {
			return all[i].Times.CreatedAt > all[j].Times.CreatedAt
		})

		json.NewEncoder(w).Encode(&all)
	}
}

func TransactionsPending(curs *currencies.Currencies, modl *model.Model) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all := modl.Reqs()
		pending := make([]otc.Request, 0)

		for _, req := range all {
			if req.Status != otc.DONE {
				pending = append(pending, req)
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
		all := modl.Reqs()
		completed := make([]otc.Request, 0)

		for _, req := range all {
			if req.Status == otc.DONE {
				completed = append(completed, req)
			}
		}

		sort.Slice(completed, func(i, j int) bool {
			return completed[i].Times.CreatedAt > completed[j].Times.CreatedAt
		})

		json.NewEncoder(w).Encode(&completed)
	}
}
