package servd

import (
	"flag"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	// . "github.com/skycoin/services/coin-api"
)

var (
	srvaddr = flag.String("srv", "localhost:12345", "RPC listener address")
)

func init() {
	flag.Parse()
}

// Start starts the server
func Start() (*echo.Echo, error) {
	// Add handlers for all currencies here
	// handlers := map[string]func(request Request) *Response{
	// 	"btc": BtcHandler,
	// }

	// Create new server
	// rpcServer := NewServer(*srvaddr, handlers)
	// Register shutdown handler
	// registerShutdownHandler(rpcServer)
	// Start server
	// rpcServer.Start()
	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.DefaultGzipConfig))
	e.Use(middleware.RecoverWithConfig(middleware.DefaultRecoverConfig))

	hMulti := newHandlerMulti()
	hBTC := newHandlerBTC()

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

	// BTC generate address, private keys, pubkeys from deterministic seed
	btcGroup.POST("/address/", hBTC.generateSeed)

	// BTC check the status of a transaction (tracks transactions by transaction hash)
	btcGroup.GET("/transaction/:transid", hBTC.checkTransaction)

	err := e.StartAutoTLS(":443")
	e.Logger.Fatal(err)

	return e, err
}

// func registerShutdownHandler(server *Server) {
// 	go func() {
// 		interruptChannel := make(chan os.Signal, 1)
// 		signal.Notify(interruptChannel, syscall.SIGINT)

// 		// Listen for initial shutdown signal and close the returned
// 		// channel to notify the caller.
// 		select {
// 		case sig := <-interruptChannel:
// 			fmt.Printf("Received signal (%s).  Shutting down...\n", sig)
// 			server.ShutDown()
// 		}

// 		// Listen for repeated signals and display a message so the user
// 		// knows the shutdown is in progress and the process is not
// 		// hung.
// 		for {
// 			select {
// 			case sig := <-interruptChannel:
// 				fmt.Printf("Received signal (%s).  Already "+
// 					"shutting down...", sig)
// 				os.Exit(1)
// 			}
// 		}
// 	}()
// }
