package capi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/skycoin/services/otc/types"
)

type Balance struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		Balance uint64 `json:"balance"`
		Address string `json:"address"`
	} `json:"result"`
}

func GetBalance(p string, c types.Currency, a types.Drop) (uint64, error) {
	client := &http.Client{Timeout: time.Second * 10}

	// prepare request
	req, err := http.NewRequest("GET", p+"/api/v1/btc/address/"+string(a), nil)
	if err != nil {
		return 0, err
	}

	// execute request
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	// unmarshal json
	var balance *Balance
	if err = json.NewDecoder(res.Body).Decode(&balance); err != nil {
		return 0, err
	}

	// check for error
	if balance.Status != "Ok" {
		return 0, fmt.Errorf(balance.Message)
	}

	return balance.Result.Balance, nil
}
