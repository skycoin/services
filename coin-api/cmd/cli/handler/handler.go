package handler

import "github.com/labstack/echo"

// Multi is multicoin handler
type Multi struct{}

// NewMulti returns multicoin handler
func NewMulti() *Multi {
	return &Multi{}
}

// GenerateKeyPair generates seed for multicurrency
func (h *Multi) GenerateKeyPair(e echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

// GenerateAddress generates seed for multicurrency
func (h *Multi) GenerateAddress(e echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

// CheckBalance checks balance for multi currency
func (h *Multi) CheckBalance(e echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

// SignTransaction signs transaction for multi currency
func (h *Multi) SignTransaction(e echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

// InjectTransaction injects transaction for any currency
func (h *Multi) InjectTransaction(e echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

// CheckTransaction checks transaction state for any currency
func (h *Multi) CheckTransaction(e echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}
