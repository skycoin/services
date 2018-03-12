package servd

type BalanceChecker interface {
	CheckBalance(string) (float64, error)
}
