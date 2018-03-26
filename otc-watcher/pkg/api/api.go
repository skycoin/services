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
			err1    error
			err2    error
		)

		if err1 = json.NewDecoder(r.Body).Decode(&req); err1 != nil {
			log.Println(err1)
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		log.Printf("Request for balance of address %s in %s\n", req.Address, req.Currency)

		if outputs, err1 = scnr.Outputs(req); err1 != nil {
			if err1 == scanner.ErrAddressMissing {
				// Register address if it missing in watch-list
				log.Printf("Register address %s", req.Address)

				if err2 = scnr.Register(req); err2 != nil {
					log.Println(err2)
					http.Error(w, err2.Error(), http.StatusInternalServerError)
					return
				}

				http.NotFound(w, r)
				return
			}

			http.Error(w, err1.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(&outputs)
	}
}
