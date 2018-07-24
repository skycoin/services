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

func (c *Custom) Update(service, version string) error {
	if c.services[service].IsLock() {
		return fmt.Errorf("service %s is already being updated", service)
	}
	c.services[service].Lock()

	customCmd, statusChan := createAndLaunch(c, service, version)
	ticker := time.NewTicker(time.Second * 2)

	go logStdout(ticker, customCmd)

	go timeoutCmd(c, service, customCmd)

	go waitForExit(statusChan, c, service)

	return nil
}

func createAndLaunch(c *Custom, service string, version string) (*cmd.Cmd, <-chan cmd.Status) {
	command := buildCommand(c, service, version)
	logrus.Info("running command: ", command)
	customCmd := cmd.NewCmd(c.services[service].ScriptInterpreter, command...)
	statusChan := customCmd.Start()
	return customCmd, statusChan
}

func buildCommand(c *Custom, service, version string) []string {
	command := []string{
		c.services[service].UpdateScript,
		c.services[service].LocalName,
		version,
	}
	return append(command, c.services[service].ScriptExtraArguments...)
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

func timeoutCmd(c *Custom, service string, customCmd *cmd.Cmd) {
	<-time.After(c.services[service].ScriptTimeout)
	customCmd.Stop()
}

func waitForExit(statusChan <-chan cmd.Status, c *Custom, service string) {
	finalStatus := <-statusChan
	logrus.Infof("%s exited with: %d", finalStatus.Cmd, finalStatus.Exit)
	c.services[service].Unlock()
}
