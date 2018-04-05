package mock

import (
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/stretchr/testify/mock"
)

// GuiClientMock represents a mock for Skycoin Web RPC API
type GuiClientMock struct {
	mock.Mock
}

// GetTransactionByID method mock
func (m *GuiClientMock) Transaction(s string) (*visor.TransactionResult, error) {
	args := m.Called(s)
	return args.Get(0).(*visor.TransactionResult), args.Error(1)
}

// InjectTransactionString method mock
func (m *GuiClientMock) InjectTransaction(s string) (string, error) {
	args := m.Called(s)
	return args.String(0), args.Error(1)
}

func (m *GuiClientMock) Balance(addresses []string) (*wallet.BalancePair, error) {
	return nil, nil
}
