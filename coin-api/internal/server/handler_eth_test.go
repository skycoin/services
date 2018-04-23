package server

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
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
		Code   int                `json:"expectedCode"`
		Result ethKeyPairResponse `json:"result"`
	}{}

	err := json.NewDecoder(resp.Body).Decode(result)

	if err != nil {
		t.Error(err)
	}

	if result.Code != resp.Code {
		t.Errorf("Wrong response expectedCode expected %d actual %d", result.Code, resp.Code)
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

func TestHandlerEthGetAddressBalance(t *testing.T) {
	testData := []struct {
		address         string
		expectedBalance int64
		err             string
	}{
		{
			"0x884348184ada7f363b2603770f03916d6137b1bf",
			10,
			"",
		},
		{
			"",
			0,
			"address not found",
		},
	}

	for _, test := range testData {
		e := echo.New()

		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		ctx := e.NewContext(req, resp)

		m := MockEthService{
			getBalance: func(ctx context.Context, addresses common.Address) (int64, error) {
				if len(test.err) > 0 {
					return 0, errors.New(test.err)
				}

				return test.expectedBalance, nil
			},
		}

		handler := HandlerEth{
			m,
		}

		ctx.SetParamNames("address")
		ctx.SetParamValues(test.address)

		handler.GetAddressBalance(ctx)

		if len(test.err) != 0 {
			result := &struct {
				Status string `json:"status"`
				Code   int    `json:"code"`
				Result string `json:"result"`
			}{}

			err := json.NewDecoder(resp.Body).Decode(result)

			if err != nil {
				t.Error(err)
				return
			}

			if test.err != result.Result {
				t.Errorf("Wrong error message expected %s actual %s",
					test.err, result.Result)
				return
			}

			return
		}

		result := &struct {
			Status string             `json:"status"`
			Code   int                `json:"expectedCode"`
			Result ethBalanceResponse `json:"result"`
		}{}

		err := json.NewDecoder(resp.Body).Decode(result)

		if err != nil {
			t.Error(err)
			return
		}

		if result.Result.Address != test.address {
			t.Errorf("Wrong address in response expected %s actual %s",
				test.address, result.Result.Address)
			return
		}

		if result.Result.Balance != test.expectedBalance {
			t.Errorf("Wrong result balance expected %d actual %d",
				test.expectedBalance, result.Result.Balance)
			return
		}
	}
}
