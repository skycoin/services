package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/src/starter"
)

var configFile string

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	flag.StringVar(&configFile,"config", "../configuration.yml", "yaml configuration file")
	flag.Parse()

	configuration := config.New(configFile)
	s := starter.New(configuration)
	s.Start()
	<- sigs
	s.Stop()
}

func stringPickNonZero(confValue, flagValue string) string {
	if confValue == "" {
		return flagValue
	}

	return confValue
}

func stringSlicePickNonZero(confValue, flagValue []string) []string {
	if confValue == nil {
		return flagValue
	}

	return confValue
}

func intPickNonZero(confValue, flagValue int) int{
	if confValue == 0 && flagValue != 0 {
		return flagValue
	}

	return confValue
}

func durationPickNonZero(confValue, flagValue time.Duration) time.Duration {
	if confValue == 0 {
		return flagValue
	}

	return confValue
}
