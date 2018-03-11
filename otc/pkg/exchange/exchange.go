package exchange

import (
	"encoding/json"
	"fmt"
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

// GetBTCValue gets the current value (in satoshis) of 1 SKY from an exchange.
// Currently, the exchange is Cryptopia. Might add more/fallback later.
func GetBTCValue() (uint64, error) {
	// request will timeout after 30 second
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	// prepare request
	req, err := http.NewRequest(
		"GET", "https://www.cryptopia.co.nz/api/GetMarket/SKY_BTC", nil,
	)
	if err != nil {
		return 0, err
	}

	// execute request
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	// get lastPrice from json body
	var market *Market
	if err = json.NewDecoder(res.Body).Decode(&market); err != nil {
		return 0, err
	}

	// check "Success" field from cryptopia and return error if needed
	if !market.Success {
		return 0, fmt.Errorf("cryptopia: %s", market.Message)
	}

	// return last price
	return uint64(market.Data.LastPrice * 10e7), nil
}
