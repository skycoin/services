package server

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/skycoin/services/coin-api/internal/eth"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	ethAddrLen       = 42
	ethPrivateKeyLen = 64
	rawAddr          = "0x51f6d925e9acfb59dfa6d3553d99f5d06b541d0c"
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

func TestEthGetAddressBalance(t *testing.T) {
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

func TestEthGetTransactionStatus(t *testing.T) {
	nonce := uint64(777)
	addr := common.HexToAddress(rawAddr)
	amount := big.NewInt(100)
	gasLimit := big.NewInt(2)
	gasPrice := big.NewInt(1)
	data := []byte("hello")

	tx := types.NewTransaction(nonce, addr, amount, gasLimit, gasPrice, data)

	chainId := big.NewInt(1)
	senderPrivKey, _ := crypto.HexToECDSA("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	signer := types.NewEIP155Signer(chainId)
	signedTx, _ := types.SignTx(tx, signer, senderPrivKey)

	testData := []struct {
		txStatus ethTxStatusResponse
		error    string
	}{
		{
			ethTxStatusResponse{
				signedTx,
				true,
			},
			"",
		},
		{
			ethTxStatusResponse{
				nil,
				true,
			},
			"Transaction not found",
		},
	}

	for _, test := range testData {
		e := echo.New()

		resp := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		ctx := e.NewContext(req, resp)

		m := MockEthService{
			getTxStatus: func(context.Context, string) (*types.Transaction, bool, error) {
				if len(test.error) > 0 {
					return nil, test.txStatus.IsPending, errors.New(test.error)
				}

				return signedTx, test.txStatus.IsPending, nil
			},
		}

		handler := HandlerEth{
			m,
		}

		handler.GetTransactionStatus(ctx)

		if len(test.error) != 0 {
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

			if test.error != result.Result {
				t.Errorf("Wrong error message expected %s actual %s",
					test.error, result.Result)
				return
			}

			return
		}

		result := &struct {
			Status string              `json:"status"`
			Code   int                 `json:"expectedCode"`
			Result ethTxStatusResponse `json:"result"`
		}{}

		json.Unmarshal(resp.Body.Bytes(), &result)

		if result.Result.IsPending != test.txStatus.IsPending {
			t.Errorf("Wrong tx status expected %t actual %t ",
				test.txStatus.IsPending, result.Result.IsPending)
			return
		}

		if result.Result.TxBody.Value().Int64() != signedTx.Value().Int64() {
			t.Errorf("Wrong tx value expected %d actual %d",
				tx.Value().Int64(), result.Result.TxBody.Value().Int64())
			return
		}
	}
}
