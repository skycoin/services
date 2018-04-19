package server

import (
	"github.com/labstack/echo"
	"github.com/skycoin/services/coin-api/internal/eth"
	"net/http"
)

type handlerEth struct {
	service *eth.EthService
}

type ethKeyPairResponse struct {
	PrivateKey string `json:"private_key"`
	Address    string `json:"address"`
}

func NewHandlerEth(nodeUrl string) (*handlerEth, error) {
	service, err := eth.NewEthService(nodeUrl)

	if err != nil {
		return nil, err
	}

	return &handlerEth{
		service: service,
	}, nil
}

func (h *handlerEth) GenerateKeyPair(ctx echo.Context) error {
	privateKey, address, err := h.service.GenerateKeyPair()

	if err != nil {
		handleError(ctx, err)
	}

	resp := ethKeyPairResponse{
		PrivateKey: privateKey,
		Address:    address,
	}

	ctx.JSONPretty(http.StatusOK, struct {
		Status string             `json:"status"`
		Code   int                `json:"code"`
		Result ethKeyPairResponse `json:"result"`
	}{
		Status: "",
		Code:   http.StatusOK,
		Result: resp,
	}, "\t")

	return nil
}

func (h *handlerEth) GetAddressBalance(ctx *echo.Context) {

}

func (h *handlerEth) GetTransactionStatus(ctx *echo.Context) {

}
