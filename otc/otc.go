package main

import (
	"fmt"
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

var (
	CURRENCIES = currencies.New()
	CONFIG     *otc.Config
)

func init() {
	var err error

	CONFIG, err = otc.NewConfig("config.toml")
	if err != nil {
		panic(err)
	}

	SKY, err := sky.New(CONFIG)
	if err != nil {
		panic(err)
	}
	CURRENCIES.Add(otc.SKY, SKY)

	BTC, err := btc.New(CONFIG)
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
	go http.ListenAndServe(CONFIG.API.Admin.Listen, admin)
	fmt.Printf("api.admin listening at %s\n", CONFIG.API.Admin.Listen)

	public := public.New(CURRENCIES, modl)
	go http.ListenAndServe(CONFIG.API.Public.Listen, public)
	fmt.Printf("api.public listening at %s\n", CONFIG.API.Public.Listen)

	<-stop
	println("stopping")
	modl.Stop()
}
