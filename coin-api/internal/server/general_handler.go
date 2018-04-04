package server

import (
	"github.com/labstack/echo"
	"github.com/skycoin/services/coin-api/internal/multi"
	"net/http"
)

type GeneralHandler struct{}

type PingResponse struct {
	Message string
}

func (g *GeneralHandler) Ping(ctx echo.Context) error {
	resp := struct {
		Code   int
		Status string
		Result PingResponse
	}{
		Code:   0,
		Status: multi.StatusOk,
		Result: PingResponse{
			Message: "Pong",
		},
	}

	ctx.JSONPretty(http.StatusOK, &resp, "\t")

	return nil
}

func (g *GeneralHandler) List(e echo.Context) error {
	return nil
}
