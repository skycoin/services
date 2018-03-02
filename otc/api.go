package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/skycoin/services/otc/types"
	"github.com/skycoin/skycoin/src/cipher"
)

type apiBindRequest struct {
	Address      string `json:"address"`
	DropCurrency string `json:"drop_currency"`
}

type apiBindResponse struct {
	DropAddress  string `json:"drop_address"`
	DropCurrency string `json:"drop_type"`
	DropValue    uint64 `json:"drop_value"`
}

func apiBind(w http.ResponseWriter, r *http.Request) {
	req := apiBindRequest{}

	// decode json
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// decode drop_currency
	currency := types.Currency(req.DropCurrency)

	// decode skycoin address
	address, err := cipher.DecodeBase58Address(req.Address)
	if err != nil {
		http.Error(w, "invalid skycoin address", http.StatusBadRequest)
		return
	}

	if DROPPER.Connections[currency] == nil {
		// currency doesn't exist
		return
	}

	// check if model is accepting new requests
	if MODEL.Paused() {
		http.Error(w, "paused", http.StatusInternalServerError)
		return
	}

	// generate drop address
	drop, err := DROPPER.Connections[currency].Generate()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		ERRS.Printf("api: %v\n", err)
		return
	}

	// create new request
	request := &types.Request{
		Address:  types.Address(address.String()),
		Currency: currency,
		Drop:     drop,
		Metadata: &types.Metadata{
			Status:    types.DEPOSIT,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
			TxId:      "",
		},
	}

	// get sky value of currency
	value, err := DROPPER.Connections[currency].Value()
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		ERRS.Printf("api: %v\n", err)
		return
	}

	// add for processing
	if err = MODEL.AddNew(request); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		ERRS.Printf("api: %v\n", err)
		return
	}

	// send json response
	if err = json.NewEncoder(w).Encode(&apiBindResponse{
		DropAddress:  string(request.Drop),
		DropCurrency: string(request.Currency),
		DropValue:    value,
	}); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		ERRS.Printf("api: %v\n", err)
		return
	}
}

type apiStatusRequest struct {
	DropAddress  string `json:"drop_address"`
	DropCurrency string `json:"drop_currency"`
}

type apiStatusResponse struct {
	Status    string `json:"status"`
	UpdatedAt int64  `json:"updated_at"`
}

func apiStatus(w http.ResponseWriter, r *http.Request) {
	req := apiStatusRequest{}

	// decode json
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// get metadata from disk
	meta, err := MODEL.GetMetadata(
		types.Drop(req.DropAddress),
		types.Currency(req.DropCurrency),
	)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		ERRS.Printf("api: %v\n", err)
		return
	}

	// send response
	json.NewEncoder(w).Encode(&apiStatusResponse{
		Status:    string(meta.Status),
		UpdatedAt: meta.UpdatedAt,
	})
}

const OTCWorking = "WORKING"
const OTCSoldOut = "SOLD_OUT"
const OTCPaused = "PAUSED"

type getConfigurationResponse struct {
	OTCStatus string `json:"otcStatus"`
	Balance   uint32 `json:"balance"`
}

func apiGetConfiguration(w http.ResponseWriter, r *http.Request) {
	status := OTCWorking
	if MODEL.Paused() {
		status = OTCPaused
	}
	// TODO other OTC statuses

	json.NewEncoder(w).Encode(&getConfigurationResponse{
		OTCStatus: status,
		Balance:   0, // TODO
	})
}
