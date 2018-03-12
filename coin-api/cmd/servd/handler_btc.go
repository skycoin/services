package servd

import (
	"crypto/rand"
	"github.com/labstack/echo"
	"github.com/skycoin/services/coin-api/internal/btc"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"net/http"
)

type keyPairResponse struct {
	Public  string `json:"public"`
	Private []byte `json:"private"`
}

type balanceResponse struct {
	Balance float64 `json:"balance"`
	Address string  `json:"address"`
}

type addressRequest struct {
	PublicKey string `json:"key"`
}

type addressResponse struct {
	Address string `json:"address"`
}

type handlerBTC struct {
	// TODO(stgleb): extract btc service to separate interface for generating key-pair and address
	btcService *btc.ServiceBtc
	checker    Checker
}

type BtcStats struct {
	NodeStatus bool   `json:"node-status"`
	NodeHost   string `json:"node-host"`
}

func newHandlerBTC(btcAddr, btcUser, btcPass string, disableTLS bool, cert []byte, blockExplorer string) (*handlerBTC, error) {
	log.Printf("Start new BTC handler with host %s user %s", btcAddr, btcUser)
	service, err := btc.NewBTCService(btcAddr, btcUser, btcPass, disableTLS, cert, blockExplorer)

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

	if err := public.Verify(); err != nil {
		return handleError(ctx, err)
	}

	resp := struct {
		Status string          `json:"status"`
		Code   int             `json:"code"`
		Result keyPairResponse `json:"result"`
	}{
		"Ok",
		http.StatusOK,
		keyPairResponse{
			Public:  public.Hex(),
			Private: private[:],
		},
	}

	// Write response with newly created key pair
	ctx.JSON(http.StatusCreated, resp)
	return nil
}

func (h *handlerBTC) generateAddress(ctx echo.Context) error {
	var req addressRequest

	if err := ctx.Bind(&req); err != nil {
		return handleError(ctx, err)
	}

	if len(req.PublicKey) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "public key is empty")
	}

	publicKey, err := cipher.PubKeyFromHex(req.PublicKey)

	if err != nil {
		return handleError(ctx, err)
	}

	address, err := btc.ServiceBtc{}.GenerateAddr(publicKey)

	if err != nil {
		return handleError(ctx, err)
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
	txId := ctx.Param("transid")
	status, err := h.checker.CheckTxStatus(txId)

	if err != nil {
		return handleError(ctx, err)
	}

	ctx.JSONPretty(http.StatusOK, struct {
		Status string        `json:"status"`
		Code   int           `json:"code"`
		Result *btc.TxStatus `json:"result"`
	}{
		Status: "",
		Code:   http.StatusOK,
		Result: status,
	}, "\t")

	return nil
}

func (h *handlerBTC) checkBalance(ctx echo.Context) error {
	address := ctx.Param("address")
	balance, err := h.checker.CheckBalance(address)

	if err != nil {
		return handleError(ctx, err)
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

func handleError(ctx echo.Context, err error) error {
	return ctx.JSONPretty(http.StatusOK, struct {
		Status string `json:"status"`
		Code   int    `json:"code"`
		Result string `json:"result"`
	}{
		Status: "",
		Code:   http.StatusOK,
		Result: err.Error(),
	}, "\t")
}
