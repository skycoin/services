package exchange

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type market struct {
	success bool `json:"Success"`
	data    struct {
		lastPrice float64 `json:"LastPrice"`
	} `json:"Data"`
}

func GetBTCPrice() (float64, error) {
	var (
		client = &http.Client{}
		m      *market
	)

	request, err := http.NewRequest("GET", "https://www.cryptopia.co.nz/api/GetMarket/SKY_BTC", nil)
	if err != nil {
		return 0.0, err
	}
	request.Header.Add("Accept-Encoding", "gzip")

	response, err := client.Do(request)
	if err != nil {
		return 0.0, err
	}
	defer response.Body.Close()

	out, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0.0, err
	}

	if err = json.Unmarshal(out, &m); err != nil {
		return 0.0, err
	}

	fmt.Println(m)

	return 0.0, nil
}
