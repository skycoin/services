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
	MODEL   *model.Model
	CONFIG  *types.Config
	DROPPER *dropper.Dropper
	ERRS    *log.Logger
)

func init() {
	var err error

	// load config file from disk
	CONFIG, err = types.NewConfig("config.toml")
	if err != nil {
		panic(err)
	}

	// error logging
	ERRS = log.New(os.Stdout, types.LOG_ERRS, types.LOG_FLAGS)

	// manages connection to btc daemon
	DROPPER, err = dropper.NewDropper(CONFIG)
	if err != nil {
		panic(err)
	}

	// manages connection and wallet for skycoin
	skycoin, err := skycoin.NewConnection(CONFIG)
	if err != nil {
		panic(err)
	}

	// actor for scanning drops
	scanner, err := scanner.NewScanner(CONFIG, DROPPER)
	if err != nil {
		panic(err)
	}
	scanner.Start()

	// actor for sending sky from otc
	sender, err := sender.NewSender(CONFIG, skycoin, DROPPER)
	if err != nil {
		panic(err)
	}
	sender.Start()

	// actor for confirming skycoin transactions
	monitor, err := monitor.NewMonitor(CONFIG, skycoin)
	if err != nil {
		panic(err)
	}
	monitor.Start()

	// actor for managing state
	MODEL, err = model.NewModel(CONFIG, scanner, sender, monitor, ERRS)
	if err != nil {
		panic(err)
	}
	MODEL.Start()
}

func main() {
	// for graceful shutdown / cleanup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		<-stop
		MODEL.Stop()
		os.Exit(0)
	}()

	// public web server
	// private web server

	http.HandleFunc("/api/bind", apiBind)
	http.HandleFunc("/api/status", apiStatus)
	http.HandleFunc("/api/config", apiGetConfigurationi)

	println("listening on " + CONFIG.Api.Listen)
	if err := http.ListenAndServe(CONFIG.Api.Listen, nil); err != nil {
		panic(err)
	}
}
