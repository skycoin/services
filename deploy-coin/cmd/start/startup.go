package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/skycoin/skycoin/src/coin"

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

func startDaemon(cfg NodeConfig, gb coin.SignedBlock) (func(), error) {
	return nil, nil
}

func startRPC(cfg NodeConfig) (func(), error) {
	return nil, nil
}

func startGUI(cfg NodeConfig) (func(), error) {
	return nil, nil
}

func makeDir(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}

	return os.Mkdir(dir, 0777)
}
