package api

import (
	"encoding/json"
	"net/http"

	"log"

	"github.com/skycoin/services/otc/pkg/otc"

	"github.com/skycoin/services/otc-watcher/pkg/scanner"
)

func New(scnr *scanner.Scanner) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/outputs", Outputs(scnr))
	return mux
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

		log.Printf("Request for balance of address %s in %s\n", req.Address, req.Currency)

		if outputs, err = scnr.Outputs(req); err != nil {
			if err == scanner.ErrAddressMissing {
				// Register address if it missing in watch-list
				log.Printf("Register address %s", req.Address)

				if err := scnr.Register(req); err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				http.NotFound(w, r)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(&outputs)
	}
}
