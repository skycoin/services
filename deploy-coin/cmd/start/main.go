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
	"sync"
	"syscall"
	"time"

	"github.com/skycoin/services/deploy-coin/common"
	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/gui"
	"github.com/skycoin/skycoin/src/util/browser"
)

var (
	// Version node version which will be set when build wallet by LDFLAGS
	Version = "0.20.0-dev"
	// Commit id
	Commit = ""
)

func main() {
	var (
		cfgPath = flag.String("config", "", "path to JSON configuration file for coin")

		runMaster = flag.Bool("master", false, "run node as master and distribute initial coin volume")
		runRPC    = flag.Bool("rpc", false, "run web RPC service")
		runGUI    = flag.Bool("gui", false, "lanuch web GUI for node in browser")

		port    = flag.Int("port", 0, "override port from config")
		rpcPort = flag.Int("rpcPort", 0, "override rpcPort from config")
		guiPort = flag.Int("guiPort", 0, "override guiPort from config")
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
	nodeCfg, err := makeNodeConfig(cfg)
	if err != nil {
		logger.Fatalf("invalid coin node configuration - %s", err)
	}

	// Override node config with command line parameters
	nodeCfg.RunMaster = *runMaster
	if nodeCfg.RunMaster {
		nodeCfg.Arbitrating = true
	}

	if *port > 0 {
		nodeCfg.Port = *port
	}

	if *rpcPort > 0 {
		nodeCfg.RPCInterfacePort = *rpcPort
	}

	if *guiPort > 0 {
		nodeCfg.WebInterfacePort = *guiPort
	}

	// Init general stuff
	closeLog, err := initLogger(nodeCfg)
	if err != nil {
		logger.Fatalf("failed to init logging - %s", err)
	}

	initPprof(nodeCfg)
	catchDebug()

	// Init daemon
	daemon, err := initDaemon(nodeCfg, cfg.Public.TrustedPeers)
	if err != nil {
		logger.Fatalf("failed to init node daemon - %s", err)
	}

	// Init web RPC
	var webRPC *webrpc.WebRPC
	if *runRPC {
		webRPC, err = initWebRPC(nodeCfg, daemon)
		if err != nil {
			logger.Fatalf("failed to init web RPC - %s", err)
		}
	}

	// Init web GUI
	var (
		webGUI  *gui.Server
		guiAddr string
	)
	if *runGUI {
		if webGUI, guiAddr, err = initWebGUI(nodeCfg, daemon); err != nil {
			logger.Fatalf("failed to start web GUI - %s", err)
		}
	}

	var (
		runWg sync.WaitGroup
		errCh = make(chan error, 10)
	)

	// Start daemon
	runWg.Add(1)
	go func() {
		defer runWg.Done()
		if err := daemon.Run(); err != nil {
			logger.Errorf("failled to run daemon - %s", err)
			errCh <- err
		}
	}()

	// Start web RPC
	if *runRPC {
		runWg.Add(1)
		go func() {
			defer runWg.Done()
			if err := webRPC.Run(); err != nil {
				logger.Errorf("failed to run web RPC - %s", err)
				errCh <- err
			}
		}()
	}

	// Start web GUI
	if *runGUI {
		runWg.Add(1)
		go func() {
			defer runWg.Done()
			if err := webGUI.Serve(); err != nil {
				logger.Errorf("failed to run web GUI - %s", err)
				errCh <- err
			}
		}()

		// Start web browser
		runWg.Add(1)
		go func() {
			defer runWg.Done()

			// Wait a moment just to make sure the http interface is up
			time.Sleep(time.Microsecond * 100)

			logger.Info("launching system browser with %s", guiAddr)
			if err := browser.Open(guiAddr); err != nil {
				logger.Errorf("failed to opend browser for web GUI - %s", err)
			}
		}()
	}

	// Distribute initial coin volume
	if nodeCfg.RunMaster {
		// Master will create distribuion transaction only in case of "emtpy" blockchain
		distTx := true
		if daemon.Visor.HeadBkSeq() > 0 {
			distTx = false
			logger.Warning("blockchain height is greater then zero - will not run coin distribution")
		}

		if distTx {
			runWg.Add(1)
			go func() {
				defer runWg.Done()

				time.Sleep(time.Second * 2)

				tx, err := makeDistributionTx(nodeCfg, cfg.Public.GenesisWallet, daemon)
				if err == nil {
					_, _, err = daemon.Visor.InjectTransaction(tx)
				}

				if err != nil {
					logger.Errorf("failed to run transaction to distribute coin volume - %s", err)
					errCh <- err
				}
			}()
		}
	}

	// Wait for SIGINT or startup error
	select {
	case <-catchInterrupt():
	case err := <-errCh:
		logger.Errorf("failed to start node -%s", err)
	}

	// Shutdown node
	logger.Info("Shutting down...")

	if *runGUI {
		webGUI.Shutdown()
	}

	if webRPC != nil {
		webRPC.Shutdown()
	}

	if daemon != nil {
		daemon.Shutdown()
	}

	runWg.Wait()
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
