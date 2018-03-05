package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/skycoin/services/coin-api/cmd/servd"
	"log"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.toml", "config file path")
	flag.Parse()
}

func main() {
	var config = &servd.Config{}
	_, err := toml.DecodeFile(configFile, config)

	if err != nil {
		log.Fatal(err)
	}

	servd.Start(config)
}
