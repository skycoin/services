package main

import (
	"github.com/labstack/echo"
)

type handlerMulti struct {
}

func newHandlerMulti() *handlerMulti {
	return &handlerMulti{}
}

func (h *handlerMulti) generateSeed(e echo.Context) error {
	return nil
}

func (h *handlerMulti) checkBalance(e echo.Context) error {
	return nil
}

func (h *handlerMulti) signTransaction(e echo.Context) error {
	return nil
}

func (h *handlerMulti) injectTransaction(e echo.Context) error {
	return nil
}

func (h *handlerMulti) checkTransaction(e echo.Context) error {
	return nil
}
