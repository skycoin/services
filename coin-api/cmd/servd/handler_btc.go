package servd

import (
	"crypto/rand"
	"net/http"

	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	"github.com/skycoin/services/coin-api/internal/btc"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
)

type keyPairResponse struct {
	Public  string `json:"public"`
	Private string `json:"private"`
}

type balanceResponse struct {
	Balance decimal.Decimal `json:"balance"`
	Address string          `json:"address"`
}

type addressRequest struct {
	PublicKey string `json:"key"`
}

type addressResponse struct {
	Address string `json:"address"`
}

type handlerBTC struct {
	btcService *btc.ServiceBtc
	checker    BalanceChecker
}

type BtcStats struct {
	NodeStatus bool   `json:"node-status"`
	NodeHost   string `json:"node-host"`
}

func newHandlerBTC(btcAddr, btcUser, btcPass string, disableTLS bool, cert []byte) (*handlerBTC, error) {
	log.Printf("Start new BTC handler with host %s user %s", btcAddr, btcUser)
	service, err := btc.NewBTCService(btcAddr, btcUser, btcPass, disableTLS, cert)

	if err != nil {
		return nil, err
	}

	return &handlerBTC{
		btcService: service,
		checker:    service,
	}, nil
}

func (h *handlerBTC) generateKeyPair(ctx echo.Context) error {
	buffer := make([]byte, 256)
	_, err := rand.Read(buffer)

	if err != nil {
		return err
	}

	public, private := btc.ServiceBtc{}.GenerateKeyPair()
	resp := struct {
		Status string          `json:"status"`
		Code   int             `json:"code"`
		Result keyPairResponse `json:"result"`
	}{
		"Ok",
		http.StatusOK,
		keyPairResponse{
			Public:  string(public[:]),
			Private: string(private[:]),
		},
	}

	// Write response with newly created key pair
	ctx.JSON(http.StatusCreated, resp)
	return nil
}

func (h *handlerBTC) generateAddress(ctx echo.Context) error {
	var req addressRequest

	if err := ctx.Bind(&req); err != nil {
		ctx.JSONPretty(http.StatusOK, &struct {
			Status string `json:"status"`
			Code   int    `json:"code"`
			Result string `json:"result"`
		}{
			Status: "",
			Code:   http.StatusOK,
			Result: err.Error(),
		}, "\t")
		return nil
	}

	if len(req.PublicKey) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "public key is empty")
	}

	publicKey, err := cipher.PubKeyFromHex(req.PublicKey)

	if err != nil {
		ctx.JSONPretty(http.StatusOK, struct {
			Status string `json:"status"`
			Code   int    `json:"code"`
			Result string `json:"result"`
		}{
			Status: "",
			Code:   http.StatusOK,
			Result: err.Error(),
		}, "\t")
		return nil
	}

	address, err := btc.ServiceBtc{}.GenerateAddr(publicKey)

	if err != nil {
		ctx.JSONPretty(http.StatusOK, struct {
			Status string `json:"status"`
			Code   int    `json:"code"`
			Result string `json:"result"`
		}{
			Status: "",
			Code:   http.StatusOK,
			Result: err.Error(),
		}, "\t")
		return nil
	}

	resp := struct {
		Status string          `json:"status"`
		Code   int             `json:"code"`
		Result addressResponse `json:"result"`
	}{
		Status: "",
		Code:   http.StatusOK,
		Result: addressResponse{
			Address: address,
		},
	}

	ctx.JSON(http.StatusCreated, resp)
	return nil
}

func (h *handlerBTC) checkTransaction(ctx echo.Context) error {
	ctx.JSONPretty(http.StatusOK, struct {
		Status string `json:"status"`
		Code   int    `json:"code"`
		Result string `json:"result"`
	}{
		Status: "",
		Code:   http.StatusOK,
		Result: "Not implemented",
	}, "\t")

	return nil
}

func (h *handlerBTC) checkBalance(ctx echo.Context) error {
	// TODO(stgleb): Check why address param is not passed
	address := ctx.ParamValues()[0]
	balance, err := h.checker.CheckBalance(address)

	if err != nil {
		ctx.JSONPretty(http.StatusOK, struct {
			Status string `json:"status"`
			Code   int    `json:"code"`
			Result string `json:"result"`
		}{
			Status: "",
			Code:   http.StatusOK,
			Result: err.Error(),
		}, "\t")
		return nil
	}

	resp := struct {
		Status string          `json:"status"`
		Code   int             `json:"code"`
		Result balanceResponse `json:"result"`
	}{
		Status: "Ok",
		Code:   http.StatusOK,
		Result: balanceResponse{
			Balance: balance,
			Address: address,
		},
	}

	ctx.JSON(http.StatusOK, resp)
	return nil
}

// Hook for collecting stats
func (h handlerBTC) CollectStatuses(stats *Status) {
	stats.Lock()
	defer stats.Unlock()
	stats.Stats["btc"] = &BtcStats{
		NodeHost:   h.btcService.GetHost(),
		NodeStatus: h.btcService.IsOpen(),
	}
}
