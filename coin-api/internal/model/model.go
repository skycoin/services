package model

// Coins represents coins
type Coins struct {
	Coins []Coin
}

// Coin represents any available coin in the wallet
type Coin struct {
	Cid      string `json:"сid"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	LastSeed string `json:"lastSeed"`
	Tm       string `json:"tm"`
	Type     string `json:"type"`
	Version  string `json:"version"`
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

// AddressResponse returns address as response
type AddressResponse struct {
	Address string `json:"address"`
}

// BalanceResponse Returns balance by given coin
type BalanceResponse struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
	Coin    Coin
}

// Transaction returns given transaction status
type Transaction struct {
	Transid string `json:"transid"`
	Status  string `json:"status"`
}

// TransactionSign represents transaction sign id
type TransactionSign struct {
	Signid string `json:"signid"`
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
