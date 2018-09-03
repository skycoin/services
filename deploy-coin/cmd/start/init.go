package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/gui"
	"github.com/skycoin/skycoin/src/util/cert"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/visor"
)

var (
	logger = logging.MustGetLogger("main")
)

func initLogger(cfg NodeConfig) (func(), error) {
	format := "[skycoin.%{module}:%{level}] %{message}"

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
	logCfg.Format = format
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

func initDaemon(cfg NodeConfig, peers []string) (*daemon.Daemon, error) {
	dc := makeDaemonConfg(cfg)

	db, err := visor.OpenDB(dc.Visor.Config.DBPath, false)
	if err != nil {
		logger.Error("Database failed to open: %v. Is another skycoin instance running?", err)
	}

	d, err := daemon.NewDaemon(dc, db, peers)
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

func initWebGUI(cfg NodeConfig, d *daemon.Daemon) (*gui.Server, string, error) {
	scheme := "http"
	if cfg.WebInterfaceHTTPS {
		scheme = "https"
	}

	host := fmt.Sprintf("%s:%d", cfg.WebInterfaceAddr, cfg.WebInterfacePort)

	// Init HTTPS certficate if necessary
	if cfg.WebInterfaceHTTPS {
		if err := cert.CreateCertIfNotExists(host,
			cfg.WebInterfaceCert, cfg.WebInterfaceKey, "Skycoind"); err != nil {
			return nil, "", fmt.Errorf("failed to create certificate for web GUI - %s", err)
		}
	}

	// Setup address
	fullAddr := fmt.Sprintf("%s://%s", scheme, host)
	logger.Critical("Full address: %s", fullAddr)
	if cfg.PrintWebInterfaceAddress {
		fmt.Println(fullAddr)
	}

	// Start web GUI
	var (
		server *gui.Server
		err    error

		guiCfg = gui.ServerConfig{
			StaticDir:   cfg.GUIDirectory,
			DisableCSRF: cfg.DisableCSRF,
		}
	)
	if cfg.WebInterfaceHTTPS {
		server, err = gui.CreateHTTPS(host, guiCfg, d, cfg.WebInterfaceCert, cfg.WebInterfaceKey)
	} else {
		server, err = gui.Create(host, guiCfg, d)
	}

	if err != nil {
		return nil, "", err
	}

	return server, fullAddr, nil
}

func makeDir(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}

	return os.Mkdir(dir, 0777)
}
