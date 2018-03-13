package main

import (
	"flag"
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/skycoin/services/coin-api/cmd/servd"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.toml", "config file path")
	flag.Parse()
}

func main() {
	f, err := os.Open(configFile)

	if err != nil {
		log.Fatal(err)
	}

	cfg := viper.New()
	cfg.SetConfigType("toml")
	cfg.AddConfigPath(".")
	cfg.ReadConfig(f)

	servd.Start(cfg)
}
