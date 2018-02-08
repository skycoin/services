package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/gui"
	"github.com/skycoin/skycoin/src/util/browser"
	"github.com/skycoin/skycoin/src/util/cert"
	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	logger = logging.MustGetLogger("main")
)

func initLogger(cfg NodeConfig) (func(), error) {
	modules := []string{
		"main",
		"daemon",
		"coin",
		"gui",
		"file",
		"visor",
		"wallet",
		"gnet",
		"pex",
		"webrpc",
	}

	logCfg := logging.DevLogConfig(modules)
	logCfg.Format = cfg.LogFmt
	logCfg.Colors = true
	logCfg.Level = "debug"

	var logFD *os.File
	if cfg.Logtofile {

		logDir := filepath.Join(cfg.DataDirectory, "logs")
		if err := makeDir(logDir); err != nil {
			return nil, fmt.Errorf("failed to create log dir, %s", err)
		}

		tf := "2006-01-02-030405"
		logFile := filepath.Join(logDir, fmt.Sprintf("%s-v%s.log", time.Now().Format(tf), Version))

		logFD, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file, %s", err)
		}

		logCfg.Output = io.MultiWriter(os.Stdout, logFD)
	}

	logCfg.InitLogger()

	closeLogFD := func() {
		logger.Info("closing log file")
		if logFD != nil {
			if err := logFD.Close(); err != nil {
				logger.Errorf("failed to close log file - %s", err)
			}
		}
	}

	return closeLogFD, nil
}

func initDaemon(cfg NodeConfig) (*daemon.Daemon, error) {
	dc := makeDaemonConfg(cfg)

	d, err := daemon.NewDaemon(dc)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func initWebRPC(cfg NodeConfig, d *daemon.Daemon) (*webrpc.WebRPC, error) {
	addr := fmt.Sprintf("%v:%v", cfg.RPCInterfaceAddr, cfg.RPCInterfacePort)
	rpc, err := webrpc.New(addr, d.Gateway)
	if err != nil {
		return nil, err
	}

	rpc.ChanBuffSize = 1000
	rpc.WorkerNum = cfg.RPCThreadNum

	return rpc, nil
}

func startWebGUI(cfg NodeConfig, d *daemon.Daemon) error {
	scheme := "http"
	if cfg.WebInterfaceHTTPS {
		scheme = "https"
	}

	host := fmt.Sprintf("%s:%d", cfg.WebInterfaceAddr, cfg.WebInterfacePort)

	// Init HTTPS certficate if necessary
	if cfg.WebInterfaceHTTPS {
		errs := cert.CreateCertIfNotExists(host,
			cfg.WebInterfaceCert, cfg.WebInterfaceKey, "Skycoind")

		if len(errs) != 0 {
			for _, err := range errs {
				logger.Error(err.Error())
			}
			return errors.New("failed to create certificate")
		}
	}

	fullAddr := fmt.Sprintf("%s://%s", scheme, host)
	logger.Critical("Full address: %s", fullAddr)

	// Start web GUI
	var err error
	if cfg.WebInterfaceHTTPS {
		err = gui.LaunchWebInterfaceHTTPS(host,
			cfg.GUIDirectory, d, cfg.WebInterfaceCert, cfg.WebInterfaceKey)
	} else {
		err = gui.LaunchWebInterface(host, cfg.GUIDirectory, d)
	}

	if err != nil {
		return err
	}

	// Start web browser
	go func() {
		// Wait a moment just to make sure the http interface is up
		time.Sleep(time.Millisecond * 100)

		logger.Info("Launching System Browser with %s", fullAddr)
		if err := browser.Open(fullAddr); err != nil {
			logger.Errorf("failed to opend boweser for web GUI - %s", err)
		}
	}()

	return nil
}

func makeDir(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}

	return os.Mkdir(dir, 0777)
}
