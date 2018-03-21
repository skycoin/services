package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"

	"github.com/skycoin/services/otc-watcher/pkg/api"
	"github.com/skycoin/services/otc-watcher/pkg/currency"
	"github.com/skycoin/services/otc-watcher/pkg/currency/btc"
	"github.com/skycoin/services/otc-watcher/pkg/scanner"
	"github.com/skycoin/services/otc/pkg/otc"
)

var (
	WALLET_NODE = flag.String(
		"wallet_node",
		"localhost:8332",
		"btcwallet rpc server",
	)

	WALLET_USER = flag.String(
		"wallet_user",
		"otc",
		"btcwallet rpc username",
	)

	WALLET_PASS = flag.String(
		"wallet_pass",
		"otc",
		"btcwallet password",
	)

	WALLET_ACCOUNT = flag.String(
		"wallet_account",
		"otc",
		"btcwallet account name",
	)

	PORT = flag.String("port", ":8080", "http api port")

	SCANNER *scanner.Scanner
)

func init() {
	flag.Parse()

	// get btc connection
	b, err := btc.New(
		*WALLET_NODE, *WALLET_ACCOUNT, *WALLET_USER, *WALLET_PASS,
	)
	if err != nil {
		panic(err)
	}

	// get scanner using btc connection
	SCANNER, err = scanner.New(
		map[otc.Currency]currency.Connection{otc.BTC: b},
	)
	if err != nil {
		panic(err)
	}

	// start listening on http port
	//
	// TODO: https
	go http.ListenAndServe(*PORT, api.New(SCANNER))
	println("listening on" + *PORT)
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	println("stopping")
	if err := SCANNER.Stop(); err != nil {
		panic(err)
	}
}
