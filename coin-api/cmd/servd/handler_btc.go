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

type addressResponse struct {
	Address string `json:"address"`
}

type handlerBTC struct {
	btcService *btc.BTCService
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

	public, private := btc.BTCService{}.GenerateKeyPair()
	resp := keyPairResponse{
		Public:  string(public[:]),
		Private: string(private[:]),
	}

	// Write response with newly created key pair
	ctx.JSON(http.StatusOK, resp)
	return nil
}

func (h *handlerBTC) checkTransaction(ctx echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "Implement me")
}

func (h *handlerBTC) checkBalance(ctx echo.Context) error {
	address := ctx.Param("address")
	balance, err := btc.BTCService{}.CheckBalance(address)

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

func (h *handlerBTC) generateAddress(ctx echo.Context) error {
	publicKeyRaw := ctx.Param("publicKeyRaw")

	if len(publicKeyRaw) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "public key is empty")
	}

	publicKey, err := cipher.PubKeyFromHex(publicKeyRaw)

	if err != nil {
		return err
	}

	address, err := btc.BTCService{}.GenerateAddr(publicKey)

	if err != nil {
		return err
	}

	resp := addressResponse{
		Address: address,
	}

	ctx.JSON(http.StatusCreated, resp)
	return nil
}
