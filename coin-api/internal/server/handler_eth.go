package server

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/labstack/echo"
	"github.com/skycoin/services/coin-api/internal/eth"
	"net/http"
)

type EthService interface {
	GenerateKeyPair() (string, string, error)
	GetBalance(context.Context, common.Address) (int64, error)
	GetTxStatus(context.Context, string) (*types.Transaction, bool, error)
}

type HandlerEth struct {
	service EthService
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
	TxBody    *types.Transaction `json:"tx_body"`
	IsPending bool               `json:"is_pending"`
}

func NewHandlerEth(nodeUrl string) (*HandlerEth, error) {
	service, err := eth.NewEthService(nodeUrl)

	if err != nil {
		return nil, err
	}

	return &HandlerEth{
		service: service,
	}, nil
}

func (h *HandlerEth) GenerateKeyPair(ctx echo.Context) error {
	privateKey, address, err := h.service.GenerateKeyPair()

	if err != nil {
		handleError(ctx, err)
		return nil
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

func (h *HandlerEth) GetAddressBalance(ctx echo.Context) error {
	addressString := ctx.Param("address")

	address := common.HexToAddress(addressString)
	balance, err := h.service.GetBalance(ctx.Request().Context(), address)

	if err != nil {
		handleError(ctx, err)
		return nil
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

func (h *HandlerEth) GetTransactionStatus(ctx echo.Context) error {
	txHash := ctx.Param("tx")

	txStatus, isPending, err := h.service.GetTxStatus(ctx.Request().Context(), txHash)

	resp := ethTxStatusResponse{
		txStatus,
		isPending,
	}

	if err != nil {
		handleError(ctx, err)
		return nil
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
