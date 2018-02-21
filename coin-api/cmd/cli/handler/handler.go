package handler

import (
	"github.com/urfave/cli"
)

// Multi is multicoin handler
type Multi struct{}

// NewMulti returns multicoin handler
func NewMulti() *Multi {
	return &Multi{}
}

// GenerateKeyPair generates seed for multicurrency
func (h *Multi) GenerateKeyPair(ctx *cli.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

// GenerateAddress generates seed for multicurrency
func (h *Multi) GenerateAddress(ctx *cli.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

// CheckBalance checks balance for multi currency
func (h *Multi) CheckBalance(ctx *cli.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

// SignTransaction signs transaction for multi currency
func (h *Multi) SignTransaction(ctx *cli.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

// InjectTransaction injects transaction for any currency
func (h *Multi) InjectTransaction(ctx *cli.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

// CheckTransaction checks transaction state for any currency
func (h *Multi) CheckTransaction(ctx *cli.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}
