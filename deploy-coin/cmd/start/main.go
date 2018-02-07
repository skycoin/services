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

	"github.com/skycoin/skycoin/src/gui"

	"github.com/skycoin/services/deploy-coin/common"
)

var (
	// Version node version which will be set when build wallet by LDFLAGS
	Version = "0.20.0-dev"
	// Commit id
	Commit = ""
)

func main() {
	var (
		cfgPath   = flag.String("config", "", "path to JSON configuration file for coin")
		runMaster = flag.Bool("isMaster", false, "run node as master")
		runGUI    = flag.Bool("gui", false, "lanuch web GUI for node in browser")
	)
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
	nodeCfg, err := makeNodeConfig(cfg, *runMaster)
	if err != nil {
		logger.Fatalf("invalid coin node configuration - %s", err)
	}
	gb, err := makeGenesisBlock(cfg)
	if err != nil {
		logger.Fatalf("invalid genesis block - %s", err)
	}

	// Init general stuff
	closeLog, err := initLogger(nodeCfg)
	if err != nil {
		logger.Fatalf("failed to init logging - %s", err)
	}

	initPprof(nodeCfg)
	catchDebug()

	// Init node
	daemon, err := initDaemon(nodeCfg, gb)
	if err != nil {
		logger.Fatalf("failed to init node daemon - %s", err)
	}

	webRPC, err := initWebRPC(nodeCfg, daemon)
	if err != nil {
		logger.Fatalf("failed to init node web RPC - %s", err)
	}

	if *runGUI {
		if err = startWebGUI(nodeCfg, daemon); err != nil {
			logger.Fatalf("failed to start web GUI - %s", err)
		}
	}

	// Start node
	errCh := make(chan error)

	go func() {
		errCh <- daemon.Run()
	}()

	go func() {
		errCh <- webRPC.Run()
	}()

	// Wait for SIGINT or startup error
	select {
	case <-catchInterrupt():
	case err := <-errCh:
		logger.Errorf("failed to start node -%s", err)
	}

	// Shutdown node
	logger.Info("Shutting down...")

	if *runGUI {
		gui.Shutdown()
	}

	if webRPC != nil {
		webRPC.Shutdown()
	}

	if daemon != nil {
		daemon.Shutdown()
	}

	closeLog()

	logger.Info("Goodbye")
}

// Catches SIGINT
func catchInterrupt() chan struct{} {
	var (
		sigCh  = make(chan os.Signal, 1)
		quitCh = make(chan struct{})
	)

	signal.Notify(sigCh, os.Interrupt)

	go func() {
		<-sigCh
		signal.Stop(sigCh)

		close(quitCh)
	}()

	return quitCh
}

// Catches SIGUSR1 and prints internal program state
func catchDebug() {
	sigchan := make(chan os.Signal, 1)
	//signal.Notify(sigchan, syscall.SIGUSR1)
	signal.Notify(sigchan, syscall.Signal(0xa)) // SIGUSR1 = Signal(0xa)

	go func() {
		for {
			select {
			case <-sigchan:
				printProgramStatus()
			}
		}
	}()
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
