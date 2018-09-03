package main

import (

"flag"
"os"
"os/signal"
"syscall"

"github.com/skycoin/services/autoupdater/config"
"github.com/skycoin/services/autoupdater/src/starter"

)

var configFile string

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	flag.StringVar(&configFile,"config", "./configuration.yml", "yaml configuration file")
	flag.Parse()

	configuration := config.New(configFile)
	s := starter.New(configuration)
	s.Start()
	<- sigs
	s.Stop()
}

