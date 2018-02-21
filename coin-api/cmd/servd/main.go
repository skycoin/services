package servd

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Start starts the server
func Start() (*echo.Echo, error) {
	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.DefaultGzipConfig))
	e.Use(middleware.RecoverWithConfig(middleware.DefaultRecoverConfig))

	e.Pre(middleware.MethodOverride())
	hMulti := newHandlerMulti()
	// TODO(stgleb): Add arguments for creating btc handler
	hBTC, err := newHandlerBTC("", "", "", false, []byte(""))

	apiGroupV1 := e.Group("/api/v1")
	multiGroup := apiGroupV1.Group("/multi/:coin")
	btcGroup := apiGroupV1.Group("/btc")

	// ping server
	apiGroupV1.GET("/ping", hMulti.generateSeed)

	// show currencies and api's list
	apiGroupV1.GET("/list", hMulti.generateSeed)

	// generate address, private keys, pubkeys from deterministic seed
	multiGroup.POST("/address/", hMulti.generateSeed)

	// check the balance (and get unspent outputs) for an address
	multiGroup.GET("/address/:address", hMulti.checkBalance)

	// sign a transaction
	multiGroup.POST("/transaction/sign/:sign", hMulti.signTransaction)

	// inject transaction into network
	multiGroup.POST("/transaction/:netid/:transid", hMulti.injectTransaction)

	// check the status of a transaction (tracks transactions by transaction hash)
	multiGroup.GET("/transaction/:transid", hMulti.checkTransaction)

	// Generate key pair
	btcGroup.POST("/keys/", hBTC.generateKeyPair)

	// BTC generate address based on public key
	btcGroup.POST("/address/:key", hBTC.generateAddress)

	// BTC check the balance (and get unspent outputs) for an address
	btcGroup.GET("/address/:address", hBTC.checkBalance)

	// BTC check the status of a transaction (tracks transactions by transaction hash)
	btcGroup.GET("/transaction/:transid", hBTC.checkTransaction)

	err = e.StartAutoTLS(":443")
	e.Logger.Fatal(err)
	return e, err
}
