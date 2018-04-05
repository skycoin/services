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

type signTransactionRequest struct {
	SignKey           string `json:"key"`
	SourceTransaction string `json:"transaction"`
}

type injectTransactionRequest struct {
	RawTransaction string `json:"transaction"`
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

func (h *handlerMulti) generateSeed(ctx echo.Context) error {
	var req addressRequest

	if err := ctx.Bind(&req); err != nil {
		return handleError(ctx, err)
	}

	if len(req.PublicKey) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "public key is empty")
	}

	addressResponse, err := h.service.GenerateAddr(req.PublicKey)

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

		return ctx.JSONPretty(http.StatusNotFound, rsp, "\t")
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

	return ctx.JSONPretty(http.StatusCreated, &rsp, "\t")
}

func (h *handlerMulti) checkBalance(e echo.Context) error {
	address := e.Param("address")
	balanceResponse, err := h.service.CheckBalance(address)

	if err != nil {
		log.Errorf("balance checking error %v", err)
		rsp := struct {
			Status string `json:"status"`
			Code   int    `json:"code"`
			Result string `json:"result"`
		}{
			Status: multi.StatusError,
			Code:   errhandler.RPCInvalidAddressOrKey,
			Result: "Address not found",
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

// Sign transaction takes secret key and ancestor transaction, build the new on and signs it.
func (h *handlerMulti) signTransaction(ctx echo.Context) error {
	var req signTransactionRequest

	if err := ctx.Bind(&req); err != nil {
		return handleError(ctx, err)
	}

	transactionSign, err := h.service.SignTransaction(req.SignKey, req.SourceTransaction)

	if err != nil {
		log.Errorf("sign transaction error %v", err)
		rsp := struct {
			Status string                         `json:"status"`
			Code   int                            `json:"code"`
			Result *multi.TransactionSignResponse `json:"result"`
		}{
			Status: multi.StatusError,
			Code:   errhandler.RPCTransactionError,
			Result: &multi.TransactionSignResponse{},
		}
		return ctx.JSONPretty(http.StatusNotFound, &rsp, "\t")
	}

	rsp := struct {
		Status string                         `json:"status"`
		Code   int                            `json:"code"`
		Result *multi.TransactionSignResponse `json:"result"`
	}{
		Status: multi.StatusOk,
		Code:   0,
		Result: transactionSign,
	}
	return ctx.JSONPretty(http.StatusOK, &rsp, "\t")
}

// Inject transaction receives hex-encoded transaction(raw transaction) to inject
func (h *handlerMulti) injectTransaction(ctx echo.Context) error {
	var req injectTransactionRequest

	if err := ctx.Bind(&req); err != nil {
		return handleError(ctx, err)
	}

	err := h.service.InjectTransaction(req.RawTransaction)

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
		return ctx.JSONPretty(http.StatusNotFound, &rsp, "\t")
	}

	rsp := struct {
		Status string `json:"status"`
		Code   int    `json:"code"`
	}{
		Status: multi.StatusOk,
		Code:   0,
	}

	return ctx.JSONPretty(http.StatusCreated, &rsp, "\t")
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
