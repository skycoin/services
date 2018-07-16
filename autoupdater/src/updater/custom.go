package updater

import (
	"context"
	"os/exec"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/services/autoupdater/config"
)

// This package implements a custom updater. This means, a script that would be launched upon
// update notify. Two arguments would always be passed to the script: Name of the service + version.

type Custom struct {
	services map[string]config.Service
}

func newCustomUpdater(c *config.Config) *Custom {
	return &Custom{
		services: c.Services,
	}
}

func (c *Custom) Update(service, version string) error {
	logrus.Warn("Update")
	ctx, cancel := context.WithTimeout(context.Background(), c.services[service].ScriptTimeout)
	defer cancel()

	command := buildCommand(c, service, version)

	err := exec.CommandContext(ctx, c.services[service].ScriptInterpreter, command...).Run()
	if err != nil {
		return err
	}
	return nil
}

func buildCommand(c *Custom, service, version string) []string {
	command := []string{
		c.services[service].UpdateScript,
		c.services[service].LocalName,
		version,
	}
	return append(command, c.services[service].ScriptExtraArguments...)
}
