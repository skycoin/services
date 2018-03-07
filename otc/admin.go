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
	Price  uint64 `json:"price"`
	Source string `json:"source"`
}

func adminPrice(w http.ResponseWriter, r *http.Request) {
	req := adminPriceRequest{}

	// decode json
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

	DROPPER.SetValue(types.BTC, req.Price)
	DROPPER.SetValueSource(source)
}

type adminStatusResponse struct {
	Price   uint64         `json:"price"`
	Updated int64          `json:"updated"`
	Source  dropper.Source `json:"source"`
	Paused  bool           `json:"paused"`
}

func adminStatus(w http.ResponseWriter, r *http.Request) {
	price, err := DROPPER.GetValue(types.BTC)
	if err != nil {
		http.Error(w, "error getting value", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&adminStatusResponse{
		Price:   price,
		Updated: DROPPER.GetUpdated(types.BTC).Unix(),
		Source:  DROPPER.GetValueSource(),
		Paused:  MODEL.Paused(),
	})
}
