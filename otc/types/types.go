package types

import "time"

const (
	BTC Currency = "BTC"
	ETH Currency = "ETH"

	DEPOSIT Status = "waiting_deposit"
	BUY     Status = "waiting_buy"
	SEND    Status = "waiting_send"
	CONFIRM Status = "waiting_confirm"
	DONE    Status = "done"
	EXPIRED Status = "expired"

	EXCHANGE_DEPOSIT  Status = "exchange_deposit"
	EXCHANGE_CONFIRM  Status = "exchange_confirm"
	EXCHANGE_TRADE    Status = "exchange_trade"
	EXCHANGE_RETURN   Status = "exchange_return"
	EXCHANGE_RETURNED Status = "exchange_returned"
)

type (
	Address  string
	Drop     string
	Currency string
	Status   string

	Metadata struct {
		Status      Status `json:"status"`
		CreatedAt   int64  `json:"created_at"`
		UpdatedAt   int64  `json:"updated_at"`
		TxId        string `json:"tx_id"`
		BuyDrop     Drop   `json:"buy_drop"`
		BuyStatus   Status `json:"buy_status"`
		BuyHashTo   string `json:"buy_hash_to"`
		BuyHashFrom string `json:"buy_hash_from"`
	}

	Request struct {
		Address  Address   `json:"address"`
		Currency Currency  `json:"currency"`
		Drop     Drop      `json:"drop"`
		Metadata *Metadata `json:"metadata"`
	}

	Work struct {
		Request *Request
		Result  chan *Result
	}

	Result struct {
		Request *Request
		Err     error
	}

	Service interface {
		Handle(*Request) chan *Result
		Start()
		Stop()
	}

	Connection interface {
		Generate() (Drop, error)
		Send(Drop, uint64) (string, error)
		Balance(Drop) (uint64, error)
		Confirmed(string) (bool, error)
		Connected() (bool, error)
		Stop() error
	}

	Connections map[Currency]Connection
)

func (w *Work) Return(err error) {
	w.Result <- &Result{w.Request, err}
}

func (m *Metadata) Update() { m.UpdatedAt = time.Now().Unix() }

func (m *Metadata) Expired(i int) bool {
	return time.Since(time.Unix(m.UpdatedAt, 0)) >
		(time.Hour * time.Duration(i))
}
