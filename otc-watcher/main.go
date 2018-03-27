package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"

	"github.com/BurntSushi/toml"
	"github.com/skycoin/services/otc-watcher/pkg/api"
	"github.com/skycoin/services/otc-watcher/pkg/currency"
	"github.com/skycoin/services/otc-watcher/pkg/currency/btc"
	"github.com/skycoin/services/otc-watcher/pkg/scanner"
	"github.com/skycoin/services/otc/pkg/otc"
)

type Config struct {
	RpcNode       string
	RpcUser       string
	RpcPass       string
	WalletAccount string
	WalletPass    string
	ListenStr     string
}

var (
	configFile = flag.String(
		"config",
		"config.toml",
		"config file",
	)

	scnr *scanner.Scanner
)

func init() {
	flag.Parse()

	config := &Config{}
	_, err := toml.DecodeFile(*configFile, config)

	if err != nil {
		panic(err)
	}

	// get btc connection
	b, err := btc.New(
		config.WalletAccount, config.WalletPass, config.RpcNode, config.RpcUser, config.RpcPass)

	if err != nil {
		panic(err)
	}

	// get scnr using btc connection
	scnr, err = scanner.New(
		map[otc.Currency]currency.Connection{otc.BTC: b},
	)

	if err != nil {
		panic(err)
	}

	// start listening on http port
	//
	// TODO: https
	go http.ListenAndServe(config.ListenStr, api.New(scnr))
	println("listening on" + config.ListenStr)
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	println("stopping")
	if err := scnr.Stop(); err != nil {
		panic(err)
	}
}
