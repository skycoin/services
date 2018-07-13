package updater

import (
	"context"
	"os/exec"
	"time"

	"github.com/skycoin/services/autoupdater/config"
)

// This package implements a custom updater. This means, a script that would be launched upon
// update notify. Two arguments would always be passed to the script: Name of the service + version.

type Custom struct {
	// /bin/bash, /bin/sh, whatever
	Interpreter string
	// path to the script
	Script string
	// extra arguments for the script
	ScriptArguments []string
	// timeout
	Timeout time.Duration
}

func newCustomUpdater(c *config.Global) *Custom {
	return &Custom{
		Interpreter:     c.Interpreter,
		Script:          c.Script,
		ScriptArguments: c.ScriptArguments,
		Timeout:         c.Timeout,
	}
}

func (c *Custom) Update(service, version string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	command := buildCommand(c, service, version)

	err := exec.CommandContext(ctx, c.Interpreter, command...).Run()
	if err != nil {
		return err
	}
	return nil
}

func buildCommand(c *Custom, service, version string) []string {
	command := []string{
		c.Script,
		service,
		version,
	}
	return append(command, c.ScriptArguments...)
}
