package otc

import (
	"time"
)

type Order struct {
	// "transaction : output index"
	Id string `json:"id"`

	// TODO: omit json
	Address string `json:"address"`
	// TODO: omit json
	Currency Currency `json:"currency"`

	// order status
	Status Status `json:"status"`
	// bitcoin amount in satoshis
	Amount uint64 `json:"amount"`
	// purchase information
	Purchase *Purchase `json:"purchase"`
	// timestamps for order
	Times *Times `json:"times"`
	// events for order
	Events []*Event `json:"events"`
}

type Purchase struct {
	// coin source
	Source string `json:"source"`
	// price information
	Price *Price `json:"price"`
	// skycoin amount received
	Amount uint64 `json:"amount"`
	// txid of skycoin transaction to user
	TxId string `json:"txid"`
}

type Price struct {
	// price source
	Source string `json:"source"`
	// price when executed (and sent)
	Executed uint64 `json:"executed"`

	// TODO: how to get this into order ?
	// price when quoted
	// Quoted uint64 `json:"quoted"`
}

type User struct {
	// skycoin address
	Address string `json:"address"`
	// affiliate code used when user was created
	Affiliate string `json:"affiliate"`
	// deposit location
	Drop *Drop `json:"drop"`
	// list of orders
	Orders []*Order `json:"orders"`
	// timestamps for user
	Times *Times `json:"times"`
}

type Currency string

const (
	BTC Currency = "BTC"
	SKY Currency = "SKY"
	ETH Currency = "ETH"
)

type Drop struct {
	Address  string   `json:"address"`
	Currency Currency `json:"currency"`
}

type Status string

const (
	DEPOSIT Status = "waiting_deposit"
	SEND    Status = "waiting_send"
	CONFIRM Status = "waiting_confirm"
	DONE    Status = "done"
)

/*
type Request struct {
	sync.Mutex

	Affiliate string `json:"affiliate"`
	Address   string `json:"address"`
	Status    Status `json:"status"`
	TxId      string `json:"txid"`
	Rate      *Rate  `json:"rate"`
	Drop      *Drop  `json:"drop"`
	Times     *Times `json:"timestamps"`
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
*/

type Times struct {
	CreatedAt   int64 `json:"created_at"`
	UpdatedAt   int64 `json:"updated_at"`
	DepositedAt int64 `json:"deposited_at,omitempty"`
	SentAt      int64 `json:"sent_at,omitempty"`
	ConfirmedAt int64 `json:"confirmed_at,omitempty"`
}

type Work struct {
	Order *Order
	Done  chan *Result
}

type Result struct {
	Finished int64 `json:"finished"`
	Err      error `json:"error"`
}

func (w *Work) Return(err error) {
	w.Done <- &Result{time.Now().UTC().Unix(), err}
}

/*
type Work struct {
	Request *Request
	Done    chan *Result
}
*/

type Event struct {
	Id       string `json:"id,omitempty"`
	Status   Status `json:"status"`
	Finished int64  `json:"finished"`
	Err      string `json:"error,omitempty"`
}

type Output struct {
	Amount    uint64   `json:"amount"`
	Addresses []string `json:"addresses"`
}

type OutputVerbose struct {
	Amount        uint64   `json:"amount"`
	Confirmations uint64   `json:"confirmations"`
	Addresses     []string `json:"addresses,omitempty"`
	Height        uint64   `json:"height,omitempty"`
}

type Transaction struct {
	Hash          string          `json:"hash"`
	Confirmations uint64          `json:"confirmations"`
	Out           map[int]*Output `json:"out"`
}

type Block struct {
	Height       uint64                  `json:"height"`
	Transactions map[string]*Transaction `json:"transactions"`
}

type Outputs map[string]map[int]*OutputVerbose

func (o Outputs) Update(hash string, index int, output *OutputVerbose) {
	if o[hash] == nil {
		o[hash] = make(map[int]*OutputVerbose, 0)
	}

	o[hash][index] = output
}
