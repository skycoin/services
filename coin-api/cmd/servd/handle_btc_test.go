package servd

import (
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"strings"
)

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

	if len(resp.Private) == 0{
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

}
