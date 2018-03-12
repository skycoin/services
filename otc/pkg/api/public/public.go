package public

import (
	"net/http"

	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/model"
)

func New(curs *currencies.Currencies, modl *model.Model) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/bind", Bind(curs, modl))
	mux.HandleFunc("/api/status", Status(curs, modl))
	mux.HandleFunc("/api/config", Config(curs, modl))
	return mux
}
