package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/skycoin/services/otc/pkg/api/admin"
	"github.com/skycoin/services/otc/pkg/api/public"
	"github.com/skycoin/services/otc/pkg/currencies"
	"github.com/skycoin/services/otc/pkg/currencies/btc"
	"github.com/skycoin/services/otc/pkg/currencies/sky"
	"github.com/skycoin/services/otc/pkg/model"
	"github.com/skycoin/services/otc/pkg/otc"
)

var CURRENCIES = currencies.New()

func init() {
	conf, err := otc.NewConfig("config.toml")
	if err != nil {
		panic(err)
	}

	SKY, err := sky.New(conf)
	if err != nil {
		panic(err)
	}
	CURRENCIES.Add(otc.SKY, SKY)

	BTC, err := btc.New(conf)
	if err != nil {
		panic(err)
	}
	CURRENCIES.Add(otc.BTC, BTC)
}

func main() {
	// for graceful shutdown / cleanup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	modl, err := model.New(CURRENCIES)
	if err != nil {
		panic(err)
	}

	admin := admin.New(CURRENCIES, modl)
	go http.ListenAndServe(":8000", admin)

	public := public.New(CURRENCIES, modl)
	println("listening")
	go http.ListenAndServe(":8080", public)

	<-stop
	println("stopping")
	modl.Stop()
}
