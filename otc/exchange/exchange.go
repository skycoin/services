package exchange

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"
)

type Market struct {
	Success bool   `json:"Success"`
	Message string `json:"Message"`
	Data    struct {
		LastPrice float64 `json:"LastPrice"`
	} `json:"Data"`
}

func GetBTCValue() (*big.Float, error) {
	// request will timeout after 1 second
	client := &http.Client{
		Timeout: time.Second * 1,
	}

	// prepare request
	req, err := http.NewRequest(
		"GET", "https://www.cryptopia.co.nz/api/GetMarket/SKY_BTC", nil,
	)
	if err != nil {
		return nil, err
	}

	// execute request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// get lastPrice from json body
	var market *Market
	if err = json.NewDecoder(res.Body).Decode(&market); err != nil {
		return nil, err
	}

	// check "Success" field from cryptopia and return error if needed
	if !market.Success {
		return nil, fmt.Errorf("cryptopia: %s", market.Message)
	}

	// return last price
	return big.NewFloat(market.Data.LastPrice), nil
}
