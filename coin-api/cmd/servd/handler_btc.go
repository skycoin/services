package servd

import (
	"github.com/labstack/echo"
)

type handlerBTC struct {
}

func newHandlerBTC() *handlerBTC {
	return &handlerBTC{}
}

func (h *handlerBTC) checkTransaction(e echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

func (h *handlerBTC) generateSeed(e echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further,
	// deal with io.Reader interface
	return nil
}
