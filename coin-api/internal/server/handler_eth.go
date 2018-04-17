package server

import (
	"github.com/labstack/echo"
	"github.com/skycoin/services/coin-api/internal/eth"
)

type handlerEth struct {
	service *eth.EthService
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

func (h *handlerEth) GenerateKeyPair(ctx *echo.Context) {

}

func (h *handlerEth) GetAddressBalance(ctx *echo.Context) {

}

func (h *handlerEth) GetTransactionStatus(ctx *echo.Context) {

}
