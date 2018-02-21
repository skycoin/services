package servd

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type balanceChecker struct {
	expected decimal.Decimal
}

func (b balanceChecker) CheckBalance(address string) (decimal.Decimal, error) {
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
	}

	var resp keyPairResponse

	err := json.Unmarshal(rec.Body.Bytes(), &resp)

	if err != nil {
		t.Error(err)
	}

	if len(resp.Public) == 0 {
		t.Errorf("Public key cannot be empty")
	}

	if len(resp.Private) == 0 {
		t.Errorf("Private key cannot be empty")
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
	}

	var resp addressResponse

	err := json.Unmarshal(rec.Body.Bytes(), &resp)

	if err != nil {
		t.Error(err)
	}

	if len(resp.Address) == 0 {
		t.Errorf("Address key cannot be empty")
	}
}

func TestCheckBalance(t *testing.T) {
	e := echo.New()
	expected := decimal.NewFromFloat(42.0)

	checker := balanceChecker{
		expected: expected,
	}

	handler := handlerBTC{
		checker: checker,
	}

	req := httptest.NewRequest(echo.GET, "/address/02a1633cafcc01ebfb6d78e39f687a1f0995c62fc95f51ead10a02ee0be551b5dc", nil)

	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	handler.checkBalance(ctx)

	var resp balanceResponse

	if rec.Code != http.StatusOK {
		t.Errorf("Wrong status code expected %d actual %d", http.StatusOK, rec.Code)
	}

	err := json.Unmarshal(rec.Body.Bytes(), &resp)

	if err != nil {
		t.Error(err)
	}

	actual, _ := resp.Balance.Float64()
	expectedFloat, _ := expected.Float64()

	if !resp.Balance.Equal(expected) {
		t.Errorf("Wrong account balance expected %f actual %f", expectedFloat, actual)
	}
}
