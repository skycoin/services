package servd

import (
	"crypto/rand"
	"github.com/labstack/echo"
	"github.com/skycoin/services/coin-api/internal/btc"
	"net/http"
)

type responseKeyPair struct {
	Public  string `json:"public"`
	Private string `json:"private"`
}

type handlerBTC struct {
}

func newHandlerBTC() *handlerBTC {
	return &handlerBTC{}
}

func (h *handlerBTC) generateKeyPair(ctx echo.Context) error {
	buffer := make([]byte, 256)
	_, err := rand.Read(buffer)

	if err != nil {
		return err
	}

	public, private := btc.BTCService{}.GenerateKeyPair()
	resp := responseKeyPair{
		Public:  string(public[:]),
		Private: string(private[:]),
	}

	// Write response with newly created key pair
	ctx.JSON(http.StatusOK, resp)
	return nil
}

func (h *handlerBTC) checkTransaction(ctx echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

func (h *handlerBTC) checkBalance(ctx echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

func (h *handlerBTC) generateAddress(ctx echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further,
	// deal with io.Reader interface

	return nil
}