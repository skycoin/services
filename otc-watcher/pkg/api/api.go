package api

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/services/otc-watcher/pkg/scanner"
	"github.com/skycoin/services/otc/pkg/otc"
)

func New(scnr *scanner.Scanner) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", Register(scnr))
	mux.HandleFunc("/outputs", Outputs(scnr))
	return mux
}

func Register(scnr *scanner.Scanner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			req *otc.Drop
			err error
		)

		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if err = scnr.Register(req); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
	}
}

func Outputs(scnr *scanner.Scanner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			outputs otc.Outputs
			req     *otc.Drop
			err     error
		)

		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if outputs, err = scnr.Outputs(req); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(&outputs)
	}
}
