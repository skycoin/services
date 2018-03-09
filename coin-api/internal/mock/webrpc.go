package mock

import (
	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/stretchr/testify/mock"
)

// WebRPCAPIMock represents a mock for Skycoin Web RPC API
type WebRPCAPIMock struct {
	mock.Mock
}

// GetTransactionByID method mock
func (m *WebRPCAPIMock) GetTransactionByID(s string) (*webrpc.TxnResult, error) {
	args := m.Called(s)
	return args.Get(0).(*webrpc.TxnResult), args.Error(1)
}

// InjectTransactionString method mock
func (m *WebRPCAPIMock) InjectTransactionString(s string) (string, error) {
	args := m.Called(s)
	return args.String(0), args.Error(1)
}
