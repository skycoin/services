package servd

type Checker interface {
	CheckBalance(string) (interface{}, error)
	CheckTxStatus(string) (interface{}, error)
}
