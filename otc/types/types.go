package types

import "time"

const (
	BTC Currency = "BTC"
	ETH Currency = "ETH"

	DEPOSIT Status = "waiting_deposit"
	SEND    Status = "waiting_send"
	CONFIRM Status = "waiting_confirm"
	DONE    Status = "done"
	EXPIRED Status = "expired"
)

type (
	Address  string
	Drop     string
	Currency string
	Status   string

	Metadata struct {
		Status    Status `json:"status"`
		Amount    uint64 `json:"amount"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`
		TxId      string `json:"tx_id"`
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
		Count() int
		Start()
		Stop()
	}

	Connection interface {
		Generate() (Drop, error)
		Balance(Drop) (uint64, error)
		Value() (uint64, error)
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
