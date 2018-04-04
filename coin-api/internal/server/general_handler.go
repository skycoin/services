package server

import (
	"github.com/labstack/echo"
	"github.com/skycoin/services/coin-api/internal/multi"
	"net/http"
	"time"
)

const (
	deterministic    = "deterministic"
	nondeterministic = "non deterministic"
)

type GeneralHandler struct{}

type PingResponse struct {
	Message string `json:"message"`
}

type ListItem struct {
	ID        string `json:"сid"`
	Name      string `json:"name"`
	TimeStamp int64  `json:"timestamp"`
	Type      string `json:"type"`
	Version   string `json:"version"`
}

type ListResponse struct {
	List []ListItem
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

func (g *GeneralHandler) List(ctx echo.Context) error {
	resp := struct {
		Code   int
		Status string
		Result ListResponse
	}{
		Code:   http.StatusOK,
		Status: multi.StatusOk,
		Result: ListResponse{
			List: []ListItem{
				{
					ID:        "BTC",
					Name:      "bitcoin",
					TimeStamp: time.Now().Unix(),
					Type:      deterministic,
					Version:   "0.1",
				},
				{
					ID:        "SKY",
					Name:      "skycoin",
					TimeStamp: time.Now().Unix(),
					Type:      nondeterministic,
					Version:   "0.1",
				},
			},
		},
	}

	ctx.JSONPretty(http.StatusOK, &resp, "\t")
	return nil
}
