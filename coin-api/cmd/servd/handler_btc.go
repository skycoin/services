package servd

import (
	"crypto/rand"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	"github.com/skycoin/services/coin-api/internal/btc"
	"github.com/skycoin/skycoin/src/cipher"
	"net/http"
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
	service, err := btc.NewBTCService(btcAddr, btcUser, btcPass, disableTLS, cert)

	if err != nil {
		return nil, err
	}

	return &handlerBTC{
		btcService: service,
	}, nil
}

func (h *handlerBTC) generateKeyPair(ctx echo.Context) error {
	buffer := make([]byte, 256)
	_, err := rand.Read(buffer)

	if err != nil {
		return err
	}

	public, private := btc.ServiceBtc{}.GenerateKeyPair()
	resp := keyPairResponse{
		Public:  string(public[:]),
		Private: string(private[:]),
	}

	// Write response with newly created key pair
	ctx.JSON(http.StatusCreated, resp)
	return nil
}

func (h *handlerBTC) generateAddress(ctx echo.Context) error {
	var req addressRequest

	if err := ctx.Bind(&req); err != nil {
		return err
	}

	if len(req.PublicKey) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "public key is empty")
	}

	publicKey, err := cipher.PubKeyFromHex(req.PublicKey)

	if err != nil {
		return err
	}

	address, err := btc.ServiceBtc{}.GenerateAddr(publicKey)

	if err != nil {
		return err
	}

	resp := addressResponse{
		Address: address,
	}

	ctx.JSON(http.StatusCreated, resp)
	return nil
}

func (h *handlerBTC) checkTransaction(ctx echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Implement me")
}

func (h *handlerBTC) checkBalance(ctx echo.Context) error {
	address := ctx.Param("address")
	balance, err := h.checker.CheckBalance(address)

	if err != nil {
		return err
	}

	resp := balanceResponse{
		Balance: balance,
		Address: address,
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
