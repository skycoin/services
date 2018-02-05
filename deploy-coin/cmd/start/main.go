package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/skycoin/services/deploy-coin/common"
)

var (
	// Version node version which will be set when build wallet by LDFLAGS
	Version = "0.20.0-dev"
	// Commit id
	Commit = ""
)

func main() {
	cfgPath := flag.String("config", "", "path to JSON configuration file for coin")
	flag.Parse()

	var (
		cfgData []byte
		err     error
	)

	// Read coin config from file or stdin
	if *cfgPath != "" {
		cfgData, err = ioutil.ReadFile(*cfgPath)
		if err != nil {
			logger.Fatalf("failed to read JSON config file - %s", err)
		}
	} else {
		cfgData, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			logger.Fatalf("failed to read JSON config from stdin")
		}
	}

	// Parse coin config
	var cfg common.Config
	if err := json.Unmarshal(cfgData, &cfg); err != nil {
		logger.Fatalf("failed ot parse JSON config - %s", err)
	}

	// Coin node config
	nodeCfg, err := makeNodeConfig(cfg)
	if err != nil {
		logger.Fatalf("invalid coin node configuration - %s", err)
	}

	gb, err := makeGenesisBlock(cfg)
	if err != nil {
		logger.Fatalf("invalid genesis block - %s", err)
	}

	// General stuff
	closeLog, err := initLogger(nodeCfg)
	if err != nil {
		logger.Fatalf("failed to init logging - %s", err)
	}

	initPprof(nodeCfg)

	quit := make(chan struct{})
	go catchInterrupt(quit)

	go catchDebug()

	// Start node
	stopDaemon, err := startDaemon(nodeCfg, gb)
	if err != nil {
		logger.Fatalf("failed to start node daemon - %s", err)
	}

	stopRPC, err := startRPC(nodeCfg)
	if err != nil {
		logger.Fatalf("failed to start node web RPC interface")
	}

	stopRPC()
	stopDaemon()
	closeLog()

	logger.Info("Goodbye")
}

// Catches SIGINT
func catchInterrupt(quit chan<- struct{}) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	<-sigchan

	signal.Stop(sigchan)
	close(quit)
}

// Catches SIGUSR1 and prints internal program state
func catchDebug() {
	sigchan := make(chan os.Signal, 1)
	//signal.Notify(sigchan, syscall.SIGUSR1)
	signal.Notify(sigchan, syscall.Signal(0xa)) // SIGUSR1 = Signal(0xa)

	for {
		select {
		case <-sigchan:
			printProgramStatus()
		}
	}
}

func initPprof(cfg NodeConfig) {
	if cfg.ProfileCPU {
		f, err := os.Create(cfg.ProfileCPUFile)
		if err != nil {
			logger.Fatalf("failed to create cpu pprof file - %s", err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if cfg.HTTPProf {
		go func() {
			fmt.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
}

func printProgramStatus() {
	fn := "goroutine.prof"
	logger.Debug("Writing goroutine profile to %s", fn)
	p := pprof.Lookup("goroutine")
	f, err := os.Create(fn)
	defer f.Close()
	if err != nil {
		logger.Error("%v", err)
		return
	}
	err = p.WriteTo(f, 2)
	if err != nil {
		logger.Error("%v", err)
		return
	}
}
