package multi

const (
	// StatusOk is an ok result in multiwallet API
	StatusOk = "ok"
	// StatusError is an error result in multiwallet API
	StatusError = "error"
	// CodeNoError means that no error occured
	CodeNoError = 0
)

// Coins represents coins
type Coins struct {
	Coins []Coin
}

// Coin represents any available coin in the wallet
type Coin struct {
	Cid      string `json:"сid"`
	Name     string `json:"name"`
	Address  string `json:"rawAddress"`
	LastSeed string `json:"lastSeed"`
	Tm       string `json:"tm"`
	Type     string `json:"type"`
	Version  string `json:"version"`
}

// Response is a typical response with status, code and result block
type Response struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Result interface{} `json:"result"`
}

// Status returns status
type Status string

// ResponsePing returns pong on a ping request
type ResponsePing struct {
	Status Status `json:"status"`
}

// KeysResponse returns generated keypair
type KeysResponse struct {
	Public  string `json:"public"`
	Private string `json:"private"`
	Status  Status `json:"status"`
}

// AddressResponse returns rawAddress as response
type AddressResponse struct {
	Address string `json:"rawAddress"`
}

// BalanceResponse Returns balance by given coin
type BalanceResponse struct {
	Address string `json:"rawAddress"`
	Balance uint64 `json:"balance"`
	Hours   uint64 `json:"hours"`
	Coin    Coin
}

// Transaction returns given transaction status
type Transaction struct {
	Transid string `json:"transid"`
	Status  string `json:"status"`
}

// TransactionSignResponse represents transaction sign id
type TransactionSignResponse struct {
	Transaction string `json:"transaction"`
}

// StdResponseMessage represents any standard message as a response for any action which doesn't return anything
type StdResponseMessage struct {
	Status  Status `json:"status"`
	Message string `json:"message"`
}

// TransactionStatus represents a status of given transaction
type TransactionStatus struct {
	Transid string `json:"transid"`
	Status  string `json:"status"`
}
