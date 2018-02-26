package servd

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"sync"
	"net/http"
)

type Status struct{
	sync.Mutex
	Stats map[string]interface{}
}

// Start starts the server
func Start() (*echo.Echo, error) {
	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.DefaultGzipConfig))
	e.Use(middleware.RecoverWithConfig(middleware.DefaultRecoverConfig))

	// e.Pre(middleware.MethodOverride())
	hMulti := newHandlerMulti()
	// TODO(stgleb): Add arguments for creating btc handler
	hBTC, err := newHandlerBTC("", "", "", false, []byte(""))

	apiGroupV1 := e.Group("/api/v1")
	skyGroup := apiGroupV1.Group("/sky")
	btcGroup := apiGroupV1.Group("/btc")

	// ping server
	apiGroupV1.GET("/ping", hMulti.generateSeed)
	// show currencies and api's list
	apiGroupV1.GET("/list", hMulti.generateSeed)
	// generate address, private keys, pubkeys from deterministic seed
	skyGroup.POST("/address", hMulti.generateSeed)
	// check the balance (and get unspent outputs) for an address
	skyGroup.GET("/address/:address", hMulti.checkBalance)
	// sign a transaction
	skyGroup.POST("/transaction/sign/:sign", hMulti.signTransaction)
	// inject transaction into network
	skyGroup.PUT("/transaction/:netid/:transid", hMulti.injectTransaction)
	// check the status of a transaction (tracks transactions by transaction hash)
	skyGroup.GET("/transaction/:transid", hMulti.checkTransaction)
	// Generate key pair
	btcGroup.POST("/keys/", hBTC.generateKeyPair)
	// // BTC generate address based on public key
	btcGroup.POST("/address/:key", hBTC.generateAddress)
	// BTC check the balance (and get unspent outputs) for an address
	btcGroup.GET("/address/:address", hBTC.checkBalance)
	// BTC check the status of a transaction (tracks transactions by transaction hash)
	btcGroup.GET("/transaction/:transid", hBTC.checkTransaction)

	statusFunc := func(ctx echo.Context) error {
		status := Status{
			Stats: make(map[string]interface{}),
		}

		// Collect statuses from handlers
		hMulti.CollectStatus(&status)
		hBTC.CollectStatuses(&status)

		ctx.JSON(http.StatusOK, status)

		return nil
	}

	// Just for basic service health checking
	e.GET("/health", func(ctx echo.Context) error {
		ctx.NoContent(http.StatusOK)
		return nil
	})

	e.GET("/status", statusFunc)
	err = e.Start(":8080")
	e.Logger.Fatal(err)
	return e, err
}
