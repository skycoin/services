package server

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

type ethBalanceResponse struct {
	Address string `json:"address"`
	Balance int64  `json:"balance"`
}

type ethTxStatusResponse struct {
	*types.Transaction
	isPending bool
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

func (h *handlerEth) GetAddressBalance(ctx echo.Context) error {
	addressString := ctx.Param("address")

	address := common.HexToAddress(addressString)
	balance, err := h.service.GetBalance(ctx.Request().Context(), address)

	if err != nil {
		handleError(ctx, err)
	}

	resp := ethBalanceResponse{
		Address: addressString,
		Balance: balance,
	}

	ctx.JSONPretty(http.StatusOK, struct {
		Status string             `json:"status"`
		Code   int                `json:"code"`
		Result ethBalanceResponse `json:"result"`
	}{
		Status: "",
		Code:   http.StatusOK,
		Result: resp,
	}, "\t")

	return nil
}

func (h *handlerEth) GetTransactionStatus(ctx echo.Context) error {
	txHash := ctx.Param("tx")

	txStatus, isPending, err := h.service.GetTxStatus(ctx.Request().Context(), txHash)

	resp := ethTxStatusResponse{
		txStatus,
		isPending,
	}

	if err != nil {
		handleError(ctx, err)
	}

	ctx.JSONPretty(http.StatusOK, struct {
		Status string              `json:"status"`
		Code   int                 `json:"code"`
		Result ethTxStatusResponse `json:"result"`
	}{
		Status: "",
		Code:   http.StatusOK,
		Result: resp,
	}, "\t")

	return nil
}
