package servd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/skycoin/services/coin-api/internal/model"
)

func TestHandlerMulti(t *testing.T) {
	e := echo.New()
	handler := handlerMulti{}

	t.Run("TestGenerateKeyPair", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/keys", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		handler.generateKeys(ctx)
		if rec.Code != http.StatusCreated {
			t.Fatalf("wrong status, expected %d  got %d", rec.Code, http.StatusCreated)
		}
		rsp := model.KeysResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &rsp)
		if err != nil {
			t.Fatal(err)
			return
		}

		if len(rsp.Private) == 0 || len(rsp.Public) == 0 {
			t.Fatal("key cannot be empty")
			return
		}

		t.Run("TestGenerateAddress", func(t *testing.T) {
			req := httptest.NewRequest(echo.POST, "/address", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			handler.generateSeed(ctx)
			if rec.Code != http.StatusCreated {
				t.Fatalf("wrong status, expected %d  got %d", rec.Code, http.StatusCreated)
			}

			rsp := &model.AddressResponse{}
			err := json.Unmarshal(rec.Body.Bytes(), rsp)
			if err != nil {
				t.Fatal(err)
				return
			}

			if len(rsp.Address) == 0 {
				t.Fatal("key cannot be empty")
				return
			}

			t.Run("checkBalance", func(t *testing.T) {
				req := httptest.NewRequest(echo.POST, fmt.Sprintf("/address/%s", rsp.Address), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				recorder := httptest.NewRecorder()
				ctx := e.NewContext(req, recorder)
				handler.checkBalance(ctx)
				rspBalance := &model.BalanceResponse{}
				err := json.Unmarshal(rec.Body.Bytes(), rspBalance)
				if err != nil {
					t.Fatalf("error unmarshalling response: %v", err)
				}
			})
		})
	})

	t.Run("signTransaction", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/address", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	})

	t.Run("injectTransaction", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/address", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	})

	t.Run("checkTransaction", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/address", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	})

}
