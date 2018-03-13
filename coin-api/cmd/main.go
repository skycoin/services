package main

import (
	"flag"
	"github.com/skycoin/services/coin-api/cmd/servd"
	"github.com/spf13/viper"
	"log"
	"os"
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
