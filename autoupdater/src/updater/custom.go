package updater

import (
	"fmt"
	"time"

	"github.com/go-cmd/cmd"
	"github.com/sirupsen/logrus"
	"github.com/skycoin/services/autoupdater/config"
)

// This package implements a custom updater. This means, a script that would be launched upon
// update notify. Two arguments would always be passed to the script: Name of the service + version.

type Custom struct {
	services map[string]*config.Service
}

func newCustomUpdater(c *config.Config) *Custom {
	return &Custom{
		services: c.Services,
	}
}

func (c *Custom) Update(service, version string) chan error {
	errCh := make(chan error)
	localService := c.services[service]

	customCmd, statusChan := createAndLaunch(localService, version)
	ticker := time.NewTicker(time.Second * 2)

	go logStdout(ticker, customCmd)

	go timeoutCmd(localService, customCmd, errCh)

	go waitForExit(statusChan, errCh)

	return errCh
}

func createAndLaunch(service *config.Service, version string) (*cmd.Cmd, <-chan cmd.Status) {
	command := buildCommand(service, version)
	logrus.Info("running command: ", command)
	customCmd := cmd.NewCmd(service.ScriptInterpreter, command...)
	statusChan := customCmd.Start()
	return customCmd, statusChan
}

func buildCommand(service *config.Service, version string) []string {
	command := []string{
		service.UpdateScript,
		service.LocalName,
		version,
	}
	return append(command, service.ScriptExtraArguments...)
}

func logStdout(ticker *time.Ticker, customCmd *cmd.Cmd) {
	var previousLastLine int

	for range ticker.C {
		status := customCmd.Status()
		currentLastLine := len(status.Stdout)

		if currentLastLine != previousLastLine {
			for _, line := range status.Stdout[previousLastLine:] {
				logrus.Infof("script stdout: %s", line)
			}
			previousLastLine = currentLastLine
		}

	}
}

func timeoutCmd( service *config.Service, customCmd *cmd.Cmd, errCh chan error) {
	<-time.After(service.ScriptTimeout)
	customCmd.Stop()
	errCh <- fmt.Errorf("update script for service %s timed out", service.OfficialName)
}

func waitForExit(statusChan <-chan cmd.Status, errCh chan error) {
	finalStatus := <-statusChan
	logrus.Infof("%s exited with: %d", finalStatus.Cmd, finalStatus.Exit)
	if finalStatus.Exit != 0 {
		errCh <- fmt.Errorf("exit with non-zero status %d", finalStatus.Exit)
	}
	errCh <- nil
}
