package main

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/services/otc/dropper"
	"github.com/skycoin/services/otc/types"
)

type adminPauseRequest struct {
	Pause bool `json:"pause"`
}

func adminPause(w http.ResponseWriter, r *http.Request) {
	req := adminPauseRequest{}

	// decode json
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Pause {
		MODEL.Pause()
	} else {
		MODEL.Unpause()
	}
}

type adminPriceRequest struct {
	Price uint64 `json:"price"`
}

func adminPrice(w http.ResponseWriter, r *http.Request) {
	req := adminPriceRequest{}

	// decode json
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	DROPPER.Currencies[types.BTC].SetValue(dropper.INTERNAL, req.Price)
}

type adminSourceRequest struct {
	Source string `json:"source"`
}

func adminSource(w http.ResponseWriter, r *http.Request) {
	req := adminSourceRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	var source dropper.Source

	// check that source option exists
	if req.Source != "exchange" && req.Source != "internal" {
		http.Error(w, "invalid price source", http.StatusBadRequest)
		return
	} else if req.Source == "exchange" {
		source = dropper.EXCHANGE
	} else if req.Source == "internal" {
		source = dropper.INTERNAL
	}

	DROPPER.Currencies[types.BTC].SetSource(source)
}

type adminStatusResponse struct {
	Prices struct {
		Internal        uint64 `json:"internal"`
		InternalUpdated int64  `json:"internal_updated"`
		Exchange        uint64 `json:"exchange"`
		ExchangeUpdated int64  `json:"exchange_updated"`
	} `json:"prices"`

	Source dropper.Source `json:"source"`
	Paused bool           `json:"paused"`
}

func adminStatus(w http.ResponseWriter, r *http.Request) {
	resp := &adminStatusResponse{
		Source: DROPPER.Currencies[types.BTC].GetSource(),
		Paused: MODEL.Paused(),
	}

	i, iu := DROPPER.Currencies[types.BTC].Sources[dropper.INTERNAL].Get()
	resp.Prices.Internal = i
	resp.Prices.InternalUpdated = iu.UTC().Unix()

	e, eu := DROPPER.Currencies[types.BTC].Sources[dropper.EXCHANGE].Get()
	resp.Prices.Exchange = e
	resp.Prices.ExchangeUpdated = eu.UTC().Unix()

	// send response
	json.NewEncoder(w).Encode(&resp)
}
