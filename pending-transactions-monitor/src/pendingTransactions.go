package pendingTransactionsMonitor

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Monitor struct {
	NodeAddress string
}

// NewMonitor creates a new instance of the NewMonitor class
func NewMonitor(nodeAddress string) *Monitor {
	return &Monitor{
		NodeAddress: nodeAddress,
	}
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Transaction represents transaction information
type Transaction struct {
	Txid string `json:"txid"`
}

// PendingTransactionsResponse represents pending transactions response
type PendingTransactionsResponse struct {
	Transaction Transaction `json:"transaction"`
	Received    time.Time   `json:"received"`
	Checked     time.Time   `json:"checked"`
}

// GetPendingTransactions returns all pending transactions
func (m Monitor) GetPendingTransactions() ([]PendingTransactionsResponse, error) {
	response, err := get(m.NodeAddress + `/pendingTxs`)
	if err != nil {
		return nil, err
	}

	transactions := []PendingTransactionsResponse{}
	if err = json.Unmarshal(response, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}
