package servd

import (
	"net/http"

	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo"
	"github.com/skycoin/skycoin/src/visor"

	"github.com/skycoin/services/coin-api/internal/locator"
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
	rsp := h.service.GenerateKeyPair()
	data, err := json.Marshal(rsp)
	if err != nil || rsp.Code != 0 {
		log.Errorf("unable to generate key, rsp code is %d, error %v", rsp.Code, err)
		return err
	}
	return e.JSONBlob(http.StatusCreated, data)
}

func (h *handlerMulti) generateSeed(e echo.Context) error {
	key := e.QueryParam("key")
	rsp, err := h.service.GenerateAddr(key)
	data, err := json.Marshal(rsp)
	if err != nil {
		log.Errorf("error encoding response %v, code %d", err, rsp.Code)
		return err
	}
	if rsp.Code != 0 {
		e.JSONBlob(http.StatusNotFound, data)
	}

	return e.JSONBlob(http.StatusCreated, data)
}

func (h *handlerMulti) checkBalance(e echo.Context) error {
	address := e.QueryParam("address")
	rsp, err := h.service.CheckBalance(address)
	if err != nil {
		log.Errorf("balance checking error %v", err)
	}
	data, err := json.Marshal(rsp)
	if err != nil {
		log.Errorf("encoding response error %v", err)
		return err
	}
	if rsp.Code != 0 {
		return e.JSONBlob(http.StatusNotFound, data)
	}
	return e.JSONBlob(http.StatusOK, data)
}

func (h *handlerMulti) signTransaction(e echo.Context) error {
	transid := e.QueryParam("signid")
	srcTrans := e.QueryParam("sourceTrans")
	rsp, err := h.service.SignTransaction(transid, srcTrans)
	if err != nil {
		log.Errorf("sign transaction error %v", err)
		return err
	}
	data, err := json.Marshal(rsp)
	if err != nil {
		log.Errorf("error encoding %v", err)
		return err
	}
	if rsp.Code != 0 {
		e.JSONBlob(http.StatusNotFound, data)
	}
	return e.JSONBlob(http.StatusOK, data)
}

func (h *handlerMulti) injectTransaction(e echo.Context) error {
	transid := e.Param("transid")
	rsp, err := h.service.InjectTransaction(transid)
	if err != nil {
		log.Errorf("inject transaction error %v", err)
		return err
	}
	data, err := json.Marshal(rsp)
	if err != nil {
		log.Errorf("responce encoding error %v", err)
		return err
	}
	if rsp.Code != 0 {
		log.Warningf("response error %d", rsp.Code)
		return e.JSONBlob(http.StatusNotFound, data)
	}
	return e.JSONBlob(http.StatusCreated, data)
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
			"Ok",
			0,
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
