package server

import (
	"log"
	"net/http"
	"os"
	"sync"

	"context"
	"net"
	"os/signal"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

type Status struct {
	sync.Mutex
	Stats map[string]interface{}
}

// Start starts the server
func Start(config *viper.Viper) error {
	e := echo.New()
	e.Use(middleware.GzipWithConfig(middleware.DefaultGzipConfig))
	e.Use(middleware.RecoverWithConfig(middleware.DefaultRecoverConfig))

	log.Printf("Create handler for skyledger based coins on %s:%d",
		config.Sub("skycoin").GetString("Host"),
		config.Sub("skycoin").GetInt("Port"))
	hMulti := newHandlerMulti(config.Sub("skycoin").GetString("Host"),
		config.Sub("skycoin").GetInt("Port"))

	hBTC, err := newHandlerBTC(
		config.Sub("bitcoin").GetString("BlockExplorer"),
		config.Sub("bitcoin").GetString("WatcherUrl"))

	generalHandler := &GeneralHandler{}

	apiGroupV1 := e.Group("/api/v1")
	skyGroup := apiGroupV1.Group("/sky")
	btcGroup := apiGroupV1.Group("/btc")

	// ping server
	apiGroupV1.GET("/ping", generalHandler.Ping)
	// show currencies and api's list
	apiGroupV1.GET("/list", generalHandler.List)

	// generate keys
	skyGroup.POST("/keys", hMulti.generateKeys)
	// generate address
	skyGroup.POST("/address/:key", hMulti.generateSeed)
	// check the balance (and get unspent outputs) for an address
	skyGroup.GET("/address/:address", hMulti.checkBalance)
	// sign a transaction
	skyGroup.POST("/transaction/sign/:sign", hMulti.signTransaction)
	// inject transaction into network
	skyGroup.PUT("/transaction/:netid/:transid", hMulti.injectTransaction)
	// check the status of a transaction (tracks transactions by transaction hash)
	skyGroup.GET("/transaction/:transid", hMulti.checkTransaction)
	// Generate key pair
	btcGroup.POST("/keys", hBTC.generateKeyPair)
	// // BTC generate address based on public key
	btcGroup.POST("/address", hBTC.generateAddress)
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
		ctx.JSONPretty(http.StatusOK, status, "\t")

		return nil
	}

	// Just for basic service health checking
	e.GET("/health", func(ctx echo.Context) error {
		ctx.NoContent(http.StatusOK)
		return nil
	})

	e.GET("/status", statusFunc)
	// Create custom listener
	listener, err := net.Listen("tcp", config.Sub("server").GetString("ListenStr"))

	if err != nil {
		return err
	}

	e.Listener = listener
	// Signal about daemon readiness
	daemon.SdNotify(false, "READY=1")

	// Custom server with configured timeouts
	server := &http.Server{
		ReadTimeout:  config.Sub("server").GetDuration("ReadTimeout") * time.Second,
		WriteTimeout: config.Sub("server").GetDuration("WriteTimeout") * time.Second,
		IdleTimeout:  config.Sub("server").GetDuration("IdleTimeout") * time.Second,
	}

	watchDog()
	interruptHandler(server)
	// Start configured server with custom listener
	err = e.StartServer(server)

	return err
}

// Handle SIGINT
func interruptHandler(server *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		ctx, cancel := context.WithTimeout(context.Background(), server.WriteTimeout)
		defer cancel()
		log.Printf("Handle SIGINT, shutdown server gracefully with timeout %f seconds", server.WriteTimeout.Seconds())
		server.Shutdown(ctx)
	}()
}

// TODO(stgleb): Do more active checks on service like http request etc.
func watchDog() {
	go func() {
		interval, err := daemon.SdWatchdogEnabled(false)
		if err != nil || interval == 0 {
			return
		}
		for {
			daemon.SdNotify(false, "WATCHDOG=1")
			time.Sleep(interval / 3)
		}
	}()
}
