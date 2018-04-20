package server_test

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/skycoin/services/coin-api/internal/eth"
	"testing"
)

type MockEthService struct{
	getBalance func(context.Context,common.Address) (int64, error)
	getTxStatus func(context.Context,string) (*types.Transaction, bool, error)

}

func (m *MockEthService) GenerateKeyPair() (string, string, error) {
	service := eth.EthService{}
	return service.GenerateKeyPair()
}

func (m *MockEthService) GetBalance(ctx context.Context,address common.Address) (int64, error) {
	return m.getBalance(ctx, address)
}

func (m *MockEthService) GetTxStatus(ctx context.Context, txid string) (*types.Transaction, bool, error) {
	return m.getTxStatus(ctx, txid)
}

