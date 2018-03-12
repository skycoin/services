package servd

import (
	"encoding/json"
	"github.com/labstack/echo"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type balanceChecker struct {
	expected float64
}

func (b balanceChecker) CheckBalance(address string) (float64, error) {
	return b.expected, nil
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
	expected := 42.0

	checker := balanceChecker{
		expected: expected,
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
		Status string          `json:"status"`
		Code   int             `json:"code"`
		Result balanceResponse `json:"result"`
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

	if math.Abs(resp.Result.Balance-expected) > 0.0001 {
		t.Errorf("Wrong account balance expected %f actual %f", expected, resp.Result.Balance)
	}
}

// TODO(stgleb): Add test for check tx status transaction
