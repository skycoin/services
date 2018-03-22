package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/skycoin/skycoin/src/visor"

	"github.com/skycoin/services/errhandler"

	"github.com/skycoin/services/coin-api/internal/locator"
	"github.com/skycoin/services/coin-api/internal/model"
	"github.com/skycoin/services/coin-api/internal/multi"
)

type handlerMulti struct {
	service *multi.SkyСoinService
}

// MultiStats some type for checking something
type MultiStats struct {
	Message string `json:"message"`
}

func newHandlerMulti(host string, port int) *handlerMulti {
	service := multi.NewSkyService(locator.NewLocatorNode(host, port))

	return &handlerMulti{
		service: service,
	}
}

// generateKeys returns a keypair
func (h *handlerMulti) generateKeys(e echo.Context) error {
	keysResponse := h.service.GenerateKeyPair()
	rsp := struct {
		Status string              `json:"status"`
		Code   int                 `json:"code"`
		Result *model.KeysResponse `json:"result"`
	}{
		Status: model.StatusOk,
		Code:   0,
		Result: keysResponse,
	}

	return e.JSONPretty(http.StatusCreated, &rsp, "\t")
}

func (h *handlerMulti) generateSeed(e echo.Context) error {
	key := e.QueryParam("key")
	addressResponse, err := h.service.GenerateAddr(key)
	if err != nil {
		log.Errorf("error encoding response %v, code", err)
		rsp := struct {
			Status string                 `json:"status"`
			Code   int                    `json:"code"`
			Result *model.AddressResponse `json:"result"`
		}{
			Status: model.StatusError,
			Code:   errhandler.RPCInvalidAddressOrKey,
			Result: &model.AddressResponse{},
		}

		return e.JSONPretty(http.StatusNotFound, rsp, "\t")
	}

	rsp := struct {
		Status string                 `json:"status"`
		Code   int                    `json:"code"`
		Result *model.AddressResponse `json:"result"`
	}{
		Status: model.StatusOk,
		Code:   0,
		Result: addressResponse,
	}

	return e.JSONPretty(http.StatusCreated, &rsp, "\t")
}

func (h *handlerMulti) checkBalance(e echo.Context) error {
	address := e.QueryParam("address")
	balanceResponse, err := h.service.CheckBalance(address)
	if err != nil {
		log.Errorf("balance checking error %v", err)
		rsp := struct {
			Status string                 `json:"status"`
			Code   int                    `json:"code"`
			Result *model.BalanceResponse `json:"result"`
		}{
			Status: model.StatusError,
			Code:   errhandler.RPCInvalidAddressOrKey,
			Result: &model.BalanceResponse{},
		}

		return e.JSONPretty(http.StatusNotFound, rsp, "\t")
	}

	rsp := struct {
		Status string                 `json:"status"`
		Code   int                    `json:"code"`
		Result *model.BalanceResponse `json:"result"`
	}{
		Status: model.StatusOk,
		Code:   0,
		Result: balanceResponse,
	}

	return e.JSONPretty(http.StatusOK, rsp, "\t")
}

func (h *handlerMulti) signTransaction(e echo.Context) error {
	transid := e.QueryParam("signid")
	srcTrans := e.QueryParam("sourceTrans")
	transactionSign, err := h.service.SignTransaction(transid, srcTrans)
	if err != nil {
		log.Errorf("sign transaction error %v", err)
		rsp := struct {
			Status string                 `json:"status"`
			Code   int                    `json:"code"`
			Result *model.TransactionSign `json:"result"`
		}{
			Status: model.StatusError,
			Code:   errhandler.RPCTransactionError,
			Result: &model.TransactionSign{},
		}
		return e.JSONPretty(http.StatusNotFound, &rsp, "\t")
	}

	rsp := struct {
		Status string                 `json:"status"`
		Code   int                    `json:"code"`
		Result *model.TransactionSign `json:"result"`
	}{
		Status: model.StatusOk,
		Code:   0,
		Result: transactionSign,
	}
	return e.JSONPretty(http.StatusOK, &rsp, "\t")
}

func (h *handlerMulti) injectTransaction(e echo.Context) error {
	transid := e.Param("transid")
	injectedTransaction, err := h.service.InjectTransaction(transid)
	if err != nil {
		log.Errorf("inject transaction error %v", err)
		rsp := struct {
			Status string             `json:"status"`
			Code   int                `json:"code"`
			Result *model.Transaction `json:"result"`
		}{
			Status: model.StatusError,
			Code:   errhandler.RPCTransactionRejected,
			Result: &model.Transaction{},
		}
		return e.JSONPretty(http.StatusNotFound, &rsp, "\t")
	}

	rsp := struct {
		Status string             `json:"status"`
		Code   int                `json:"code"`
		Result *model.Transaction `json:"result"`
	}{
		Status: model.StatusOk,
		Code:   0,
		Result: injectedTransaction,
	}

	return e.JSONPretty(http.StatusCreated, &rsp, "\t")
}

func (h *handlerMulti) checkTransaction(ctx echo.Context) error {
	txID := ctx.Param("transid")
	status, err := h.service.CheckTransactionStatus(txID)

	if err != nil {
		ctx.JSONPretty(http.StatusOK, struct {
			Status string `json:"status"`
			Code   int    `json:"code"`
			Result string `json:"result"`
		}{
			model.StatusOk,
			0,
			err.Error(),
		}, "\t")
	}

	ctx.JSONPretty(http.StatusOK, struct {
		Status string                  `json:"status"`
		Code   int                     `json:"code"`
		Result visor.TransactionStatus `json:"result"`
	}{
		model.StatusOk,
		0,
		*status,
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
