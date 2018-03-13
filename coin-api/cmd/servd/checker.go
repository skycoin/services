package servd

import "github.com/skycoin/services/coin-api/internal/btc"

type Checker interface {
	CheckBalance(string) (float64, error)
	CheckTxStatus(string) (*btc.TxStatus, error)
}
