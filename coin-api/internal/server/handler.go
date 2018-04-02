package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/skycoin/skycoin/src/visor"

	"github.com/skycoin/services/errhandler"

	"github.com/skycoin/services/coin-api/internal/multi"
)

type handlerMulti struct {
	// TODO(stgleb): Create an interface for abstracting service.
	service *multi.SkyСoinService
}

// MultiStats some type for checking something
type MultiStats struct {
	Message string `json:"message"`
}

func newHandlerMulti(host string, port int) *handlerMulti {
	service := multi.NewSkyService(multi.NewLocatorNode(host, port))

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
		Result *multi.KeysResponse `json:"result"`
	}{
		Status: multi.StatusOk,
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
			Result *multi.AddressResponse `json:"result"`
		}{
			Status: multi.StatusError,
			Code:   errhandler.RPCInvalidAddressOrKey,
			Result: &multi.AddressResponse{},
		}

		return e.JSONPretty(http.StatusNotFound, rsp, "\t")
	}

	rsp := struct {
		Status string                 `json:"status"`
		Code   int                    `json:"code"`
		Result *multi.AddressResponse `json:"result"`
	}{
		Status: multi.StatusOk,
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
			Result *multi.BalanceResponse `json:"result"`
		}{
			Status: multi.StatusError,
			Code:   errhandler.RPCInvalidAddressOrKey,
			Result: &multi.BalanceResponse{},
		}

		return e.JSONPretty(http.StatusNotFound, rsp, "\t")
	}

	rsp := struct {
		Status string                 `json:"status"`
		Code   int                    `json:"code"`
		Result *multi.BalanceResponse `json:"result"`
	}{
		Status: multi.StatusOk,
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
			Result *multi.TransactionSign `json:"result"`
		}{
			Status: multi.StatusError,
			Code:   errhandler.RPCTransactionError,
			Result: &multi.TransactionSign{},
		}
		return e.JSONPretty(http.StatusNotFound, &rsp, "\t")
	}

	rsp := struct {
		Status string                 `json:"status"`
		Code   int                    `json:"code"`
		Result *multi.TransactionSign `json:"result"`
	}{
		Status: multi.StatusOk,
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
			Result *multi.Transaction `json:"result"`
		}{
			Status: multi.StatusError,
			Code:   errhandler.RPCTransactionRejected,
			Result: &multi.Transaction{},
		}
		return e.JSONPretty(http.StatusNotFound, &rsp, "\t")
	}

	rsp := struct {
		Status string             `json:"status"`
		Code   int                `json:"code"`
		Result *multi.Transaction `json:"result"`
	}{
		Status: multi.StatusOk,
		Code:   0,
		Result: injectedTransaction,
	}

	return e.JSONPretty(http.StatusCreated, &rsp, "\t")
}

func (h *handlerMulti) checkTransaction(ctx echo.Context) error {
	txID := ctx.QueryParam("transid")
	status, err := h.service.CheckTransactionStatus(txID)

	if err != nil {
		ctx.JSONPretty(http.StatusOK, struct {
			Status string `json:"status"`
			Code   int    `json:"code"`
			Result string `json:"result"`
		}{
			multi.StatusOk,
			0,
			err.Error(),
		}, "\t")

		return err
	}

	ctx.JSONPretty(http.StatusOK, struct {
		Status string                  `json:"status"`
		Code   int                     `json:"code"`
		Result visor.TransactionStatus `json:"result"`
	}{
		multi.StatusOk,
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
