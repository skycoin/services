package main

import (
	"flag"
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
	SKYCOIN *skycoin.Connection
	ERRS    *log.Logger

	CONFIG_PATH = flag.String(
		"config",
		"config.toml",
		"path to config .toml file",
	)
)

func init() {
	flag.Parse()
	var err error

	// load config file from disk
	CONFIG, err = types.NewConfig(*CONFIG_PATH)
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
	SKYCOIN, err = skycoin.NewConnection(CONFIG)
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
	sender, err := sender.NewSender(CONFIG, SKYCOIN, DROPPER)
	if err != nil {
		panic(err)
	}
	sender.Start()

	// actor for confirming skycoin transactions
	monitor, err := monitor.NewMonitor(CONFIG, SKYCOIN)
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

	// public facing otc api
	api := http.NewServeMux()
	api.HandleFunc("/api/bind", apiBind)
	api.HandleFunc("/api/status", apiStatus)
	api.HandleFunc("/api/config", apiGetConfiguration)
	go http.ListenAndServe(CONFIG.Api.Listen, api)

	// private facing admin api
	admin := http.NewServeMux()
	admin.HandleFunc("/api/status", adminStatus)
	admin.HandleFunc("/api/pause", adminPause)
	admin.HandleFunc("/api/price", adminPrice)
	go http.ListenAndServe(CONFIG.Admin.Listen, admin)

	println("api listening on " + CONFIG.Api.Listen)
	println("admin listening on " + CONFIG.Admin.Listen)

	<-stop
	println("stopping")
	MODEL.Stop()
}
