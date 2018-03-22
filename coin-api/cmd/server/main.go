package main

import (
	"flag"
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/skycoin/services/coin-api/internal/server"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "init", "init.toml", "init file path")
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

	log.Fatal(server.Start(cfg))
}
