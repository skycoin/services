package servd

import (
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	e := echo.New()
	handler := handlerBTC{}
	req := httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	handler.generateKeyPair(ctx)

	if rec.Code != http.StatusCreated {
		t.Errorf("Wrong http status expected %d actual %d", rec.Code, http.StatusCreated)
	}
}

func TestGenerateAddress(t *testing.T) {

}

func TestCheckBalance(t *testing.T) {

}
