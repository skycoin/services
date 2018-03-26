package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"

	"github.com/skycoin/services/coin-api/internal/btc"
)

type checker struct {
	expected *btc.BalanceResponse
	txStatus *btc.TxStatus
}

func (b checker) CheckBalance(address string) (interface{}, error) {
	return b.expected, nil
}

func (b checker) CheckTxStatus(txId string) (interface{}, error) {
	return b.txStatus, nil
}

func TestGenerateKeyPair(t *testing.T) {
	e := echo.New()
	handler := handlerBTC{}
	req := httptest.NewRequest(echo.POST, "/address", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	handler.generateKeyPair(ctx)

	if rec.Code != http.StatusCreated {
		t.Errorf("Wrong http status expected %d actual %d", rec.Code, http.StatusCreated)
		return
	}

	type response struct {
		Status string          `json:"status"`
		Code   int             `json:"code"`
		Result keyPairResponse `json:"result"`
	}

	var resp response

	err := json.Unmarshal(rec.Body.Bytes(), &resp)

	if err != nil {
		t.Error(err)
		return
	}

	if len(resp.Result.Public) == 0 {
		t.Errorf("Public key cannot be empty")
		return
	}

	if len(resp.Result.Private) == 0 {
		t.Errorf("Private key cannot be empty")
		return
	}
}

func TestGenerateAddress(t *testing.T) {
	e := echo.New()
	handler := handlerBTC{}
	body := `{"key":"02a1633cafcc01ebfb6d78e39f687a1f0995c62fc95f51ead10a02ee0be551b5dc"}`

	req := httptest.NewRequest(echo.POST, "/address", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	handler.generateAddress(ctx)

	if rec.Code != http.StatusCreated {
		t.Errorf("Wrong http status expected %d actual %d", http.StatusCreated, rec.Code)
		return
	}

	type response struct {
		Status string          `json:"status"`
		Code   int             `json:"code"`
		Result addressResponse `json:"result"`
	}

	var resp response

	err := json.Unmarshal(rec.Body.Bytes(), &resp)

	if err != nil {
		t.Error(err)
		return
	}

	if len(resp.Result.Address) == 0 {
		t.Errorf("Address key cannot be empty")
		return
	}
}

func TestCheckBalance(t *testing.T) {
	e := echo.New()
	expectedBalance := int64(42)
	balanceResp := &btc.BalanceResponse{
		Balance: expectedBalance,
	}

	checker := checker{
		expected: balanceResp,
	}

	handler := handlerBTC{
		checker: checker,
	}

	req := httptest.NewRequest(echo.GET, "/", nil)

	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	ctx.SetParamNames("address")
	ctx.SetParamValues("1M3GipkG2YyHPDMPewqTpup83jitXvBg9N")

	handler.checkBalance(ctx)

	type response struct {
		Status string              `json:"status"`
		Code   int                 `json:"code"`
		Result btc.BalanceResponse `json:"result"`
	}

	var resp response

	if rec.Code != http.StatusOK {
		t.Errorf("Wrong status code expected %d actual %d", http.StatusOK, rec.Code)
		return
	}

	err := json.Unmarshal(rec.Body.Bytes(), &resp)

	if err != nil {
		t.Error(err)
		return
	}

	if resp.Result.Balance != expectedBalance {
		t.Errorf("Wrong account balance expected %f actual %f", expectedBalance, resp.Result.Balance)
	}
}

func TestCheckTransaction(t *testing.T) {
	e := echo.New()
	expectedConfirmations := int64(42)
	timeConfirmed := time.Now().Unix()
	hash := "89f04c437ee192a28c59470c010359c50239e28df903e44778286fb56b8e6e6f"

	checker := checker{
		txStatus: &btc.TxStatus{
			Hash:          hash,
			Confirmations: expectedConfirmations,
			Confirmed:     timeConfirmed,
		},
	}

	handler := handlerBTC{
		checker: checker,
	}

	req := httptest.NewRequest(echo.GET, "/", nil)

	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	ctx.SetParamNames("transid")
	ctx.SetParamValues("89f04c437ee192a28c59470c010359c50239e28df903e44778286fb56b8e6e6f")

	handler.checkTransaction(ctx)

	type response struct {
		Status string        `json:"status"`
		Code   int           `json:"code"`
		Result *btc.TxStatus `json:"result"`
	}

	var resp response

	if rec.Code != http.StatusOK {
		t.Errorf("Wrong status code expected %d actual %d", http.StatusOK, rec.Code)
		return
	}

	err := json.Unmarshal(rec.Body.Bytes(), &resp)

	if err != nil {
		t.Error(err)
		return
	}

	if resp.Result.Confirmations != expectedConfirmations {
		t.Errorf("Wrong confirmations count expected %d actual %d", expectedConfirmations, resp.Result.Confirmations)
	}

	if resp.Result.Hash != hash {
		t.Errorf("Wrong hash value expected %s actual %s", hash, resp.Result.Hash)
	}
}
