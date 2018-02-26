package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/skycoin/services/otc/dropper"
	"github.com/skycoin/services/otc/model"
	"github.com/skycoin/services/otc/monitor"
	"github.com/skycoin/services/otc/scanner"
	"github.com/skycoin/services/otc/sender"
	"github.com/skycoin/services/otc/skycoin"
	"github.com/skycoin/services/otc/types"
)

var (
	DROPPER *dropper.Dropper
	SKYCOIN *skycoin.Connection
	SCANNER *scanner.Scanner
	SENDER  *sender.Sender
	MONITOR *monitor.Monitor
	MODEL   *model.Model

	ERRS *log.Logger
)

func main() {
	// for graceful shutdown / cleanup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	CONFIG, err := types.NewConfig("config.toml")
	if err != nil {
		panic(err)
	}

	ERRS = log.New(os.Stdout, types.LOG_ERRS, types.LOG_FLAGS)

	DROPPER, err = dropper.NewDropper(CONFIG)
	if err != nil {
		panic(err)
	}

	SKYCOIN, err = skycoin.NewConnection(CONFIG)
	if err != nil {
		panic(err)
	}

	SCANNER, err = scanner.NewScanner(CONFIG, DROPPER)
	if err != nil {
		panic(err)
	}
	SCANNER.Start()

	SENDER, err = sender.NewSender(CONFIG, SKYCOIN, DROPPER)
	if err != nil {
		panic(err)
	}
	SENDER.Start()

	MONITOR, err = monitor.NewMonitor(CONFIG, SKYCOIN)
	if err != nil {
		panic(err)
	}
	MONITOR.Start()

	MODEL, err = model.NewModel(CONFIG, SCANNER, SENDER, MONITOR, ERRS)
	if err != nil {
		panic(err)
	}
	MODEL.Start()

	go func() {
		<-stop
		println("stopping")
		MODEL.Stop()
		os.Exit(0)
	}()

	http.HandleFunc("/api/bind", apiBind)
	http.HandleFunc("/api/status", apiStatus)

	println("listening on :8080")
	if err = http.ListenAndServe(CONFIG.Api.Listen, nil); err != nil {
		panic(err)
	}
}
