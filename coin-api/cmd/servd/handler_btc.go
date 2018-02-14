package main

import (
	"github.com/labstack/echo"
)

type handlerBTC struct {
}

func newHandlerBTC() *handlerBTC {
	return &handlerBTC{}
}

func (h *handlerBTC) checkTransaction(e echo.Context) error {
	return nil
}

func (h *handlerBTC) generateSeed(e echo.Context) error {
	return nil
}
