package servd

import "github.com/shopspring/decimal"

type BalanceChecker interface {
	CheckBalance(string) (decimal.Decimal, error)
}
