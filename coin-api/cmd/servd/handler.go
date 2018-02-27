package servd

import (
	"github.com/labstack/echo"
	"github.com/skycoin/services/coin-api/internal/multi"
	"github.com/skycoin/skycoin/src/visor"
	"net/http"
)

type handlerMulti struct {
	service *multi.GenericСoinService
}

type MultiStats struct {
	Message string `json:"message"`
}

func newHandlerMulti(nodeAddr string) *handlerMulti {
	service := multi.NewMultiCoinService(nodeAddr)

	return &handlerMulti{
		service: service,
	}
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

func (h *handlerMulti) checkTransaction(ctx echo.Context) error {
	txId := ctx.Param("transid")
	status, err := h.service.CheckTransactionStatus(txId)

	if err != nil {
		ctx.JSONPretty(http.StatusOK, struct {
			Status string `json:"status"`
			Code   int    `json:"code"`
			Result string `json:"result"`
		}{
			"Ok",
			http.StatusOK,
			err.Error(),
		}, "\t")
	}

	ctx.JSONPretty(http.StatusOK, struct {
		Status string                  `json:"status"`
		Code   int                     `json:"code"`
		Result visor.TransactionStatus `json:"result"`
	}{
		"Ok",
		http.StatusOK,
		status,
	}, "\t")

	return nil
}

func (h *handlerMulti) CollectStatus(status *Status) {
	status.Lock()
	defer status.Unlock()
	status.Stats["multi"] = &MultiStats{
		Message: "Implement me",
	}
}
