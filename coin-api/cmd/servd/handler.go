package servd

import (
	"github.com/labstack/echo"
)

type handlerMulti struct {
}

type MultiStats struct {
	Message string `json:"message"`
}

func newHandlerMulti() *handlerMulti {
	return &handlerMulti{}
}

func (h *handlerMulti) generateSeed(e echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

func (h *handlerMulti) checkBalance(e echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

func (h *handlerMulti) signTransaction(e echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

func (h *handlerMulti) injectTransaction(e echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

func (h *handlerMulti) checkTransaction(e echo.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	return nil
}

func (h handlerMulti) CollectStatus(status *Status) {
	status.Lock()
	defer status.Unlock()
	status.Stats["multi"] = &MultiStats{
		Message: "Implement me",
	}
}
