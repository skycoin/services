package server

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/echo"
	"github.com/skycoin/services/coin-api/internal/eth"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	ethAddrLen       = 42
	ethPrivateKeyLen = 64
)

type MockEthService struct {
	getBalance  func(context.Context, common.Address) (int64, error)
	getTxStatus func(context.Context, string) (*types.Transaction, bool, error)
}

func (m MockEthService) GenerateKeyPair() (string, string, error) {
	service := eth.EthService{}
	return service.GenerateKeyPair()
}

func (m MockEthService) GetBalance(ctx context.Context, address common.Address) (int64, error) {
	return m.getBalance(ctx, address)
}

func (m MockEthService) GetTxStatus(ctx context.Context, txid string) (*types.Transaction, bool, error) {
	return m.getTxStatus(ctx, txid)
}

func TestGenerateEthKeyPair(t *testing.T) {
	e := echo.New()
	m := MockEthService{}

	handler := HandlerEth{
		m,
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ctx := e.NewContext(req, resp)
	handler.GenerateKeyPair(ctx)

	result := &struct {
		Status string             `json:"status"`
		Code   int                `json:"code"`
		Result ethKeyPairResponse `json:"result"`
	}{}

	err := json.NewDecoder(resp.Body).Decode(result)

	if err != nil {
		t.Error(err)
	}

	if result.Code != resp.Code {
		t.Errorf("Wrong response code expected %d actual %d", result.Code, resp.Code)
	}

	if len(result.Result.Address) != ethAddrLen {
		t.Errorf("Wrong address len expected %d actual %d",
			ethAddrLen,
			len(result.Result.Address))
	}

	if len(result.Result.PrivateKey) != ethPrivateKeyLen {
		t.Errorf("Wrong private key len expected %d actual %d",
			ethPrivateKeyLen,
			len(result.Result.PrivateKey))
	}
}
