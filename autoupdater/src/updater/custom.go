package updater

import (
	"fmt"
	"time"

	"github.com/go-cmd/cmd"
	"github.com/sirupsen/logrus"
	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/src/logger"
)

// This package implements a custom updater. This means, a script that would be launched upon
// update notify. Two arguments would always be passed to the script: Name of the service + version.

var defaultScriptTimeout = time.Minute*10

type Custom struct {
	services map[string]customServiceConfig
}

type customServiceConfig struct {
	officialName string
	localName string
	scriptInterpreter string
	scriptExtraArguments []string
	updateScript string
	tag string
	scriptTimeout time.Duration
}

func newCustomUpdater(services map[string]config.ServiceConfig) *Custom {
	customServices := make(map[string]customServiceConfig)
	for officialName, c:= range services{
		duration, err := time.ParseDuration(c.ScriptTimeout)
		if err != nil {
			duration = defaultScriptTimeout
			logrus.Warnf("cannot parse timeout duration %s of service %s configuration." +
				" setting default timeout %s", c.ScriptTimeout, c.OfficialName, duration.String())
		}

		customServices[officialName] = customServiceConfig{
			officialName: officialName,
			localName: c.LocalName,
			scriptExtraArguments: c.ScriptExtraArguments,
			scriptInterpreter: c.ScriptInterpreter,
			scriptTimeout: duration,
			tag: c.CheckTag,
			updateScript: c.UpdateScript,
		}
	}

	return &Custom{
		services: customServices,
	}
}

func (c *Custom) Update(service, version string, log *logger.Logger) chan error {
	errCh := make(chan error)
	localService := c.services[service]

	customCmd, statusChan := createAndLaunch(localService, version, log)
	ticker := time.NewTicker(time.Second * 2)

	go logStdout(ticker, customCmd, log)

	go timeoutCmd(localService, customCmd, errCh)

	go waitForExit(statusChan, errCh, log)

	return errCh
}

func createAndLaunch(service customServiceConfig, version string, log *logger.Logger) (*cmd.Cmd, <-chan cmd.Status) {
	command := buildCommand(service, version)
	log.Info("running command: ", command)
	customCmd := cmd.NewCmd(service.scriptInterpreter, command...)
	statusChan := customCmd.Start()
	return customCmd, statusChan
}

func buildCommand(service customServiceConfig, version string) []string {
	command := []string{
		service.updateScript,
		service.localName,
		version,
	}
	return append(command, service.scriptExtraArguments...)
}

func logStdout(ticker *time.Ticker, customCmd *cmd.Cmd, log *logger.Logger) {
	var previousLastLine int

	for range ticker.C {
		status := customCmd.Status()
		currentLastLine := len(status.Stdout)

		if currentLastLine != previousLastLine {
			for _, line := range status.Stdout[previousLastLine:] {
				log.Infof("script stdout: %s", line)
			}
			previousLastLine = currentLastLine
		}

	}
}

func timeoutCmd( service customServiceConfig, customCmd *cmd.Cmd, errCh chan error) {
	<-time.After(service.scriptTimeout)
	customCmd.Stop()
	errCh <- fmt.Errorf("update script for service %s timed out", service.officialName)
}

func waitForExit(statusChan <-chan cmd.Status, errCh chan error, log *logger.Logger) {
	finalStatus := <-statusChan
	log.Infof("%s exit with: %d", finalStatus.Cmd, finalStatus.Exit)
	if finalStatus.Exit != 0 {
		errCh <- fmt.Errorf("exit with non-zero status %d", finalStatus.Exit)
	}
	errCh <- nil
}
