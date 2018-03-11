package otc

import (
	"sync"
	"time"
)

type Currency string

const (
	BTC Currency = "BTC"
	SKY Currency = "SKY"
	ETH Currency = "ETH"
)

type Drop struct {
	Address  string   `json:"address"`
	Currency Currency `json:"currency"`
	Amount   uint64   `json:"amount"`
}

type Status string

const (
	NEW     Status = "new"
	DEPOSIT Status = "waiting_deposit"
	SEND    Status = "waiting_send"
	CONFIRM Status = "waiting_confirm"
	DONE    Status = "done"
)

type Request struct {
	sync.Mutex

	Address string `json:"address"`
	Status  Status `json:"status"`
	TxId    string `json:"txid"`
	Rate    *Rate  `json:"rate"`
	Drop    *Drop  `json:"drop"`
	Times   *Times `json:"timestamps"`
}

func (r *Request) Id() string {
	return r.Address + ":" + string(r.Drop.Currency) + ":" + r.Drop.Address
}

func (r *Request) Iden() string {
	return string(r.Drop.Currency) + ":" + r.Drop.Address
}

type Rate struct {
	Value  uint64 `json:"value"`
	Source string `json:"source"`
}

type Times struct {
	CreatedAt   int64 `json:"created_at"`
	DepositedAt int64 `json:"deposited_at"`
	SentAt      int64 `json:"sent_at"`
	ConfirmedAt int64 `json:"confirmed_at"`
	UpdatedAt   int64 `json:"updated_at"`
}

type Work struct {
	Request *Request
	Done    chan *Result
}

type Event struct {
	Id       string `json:"id"`
	Status   Status `json:"status"`
	Finished int64  `json:"finished"`
	Err      string `json:"error"`
}

type Result struct {
	Finished int64 `json:"finished"`
	Err      error `json:"error"`
}

func (w *Work) Return(err error) {
	w.Done <- &Result{time.Now().UTC().Unix(), err}
}
