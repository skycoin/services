package updater_test

import (
	"testing"
	"time"

	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/src/updater"
	"github.com/stretchr/testify/assert"
)

const TEST_SCRIPT = `
#!/bin/bash

echo "service {$1}"
echo "version {$2}"
shift 2

echo "arguments {$@}"
`

func TestCustom(t *testing.T) {
	customConfig := &config.Config{
		Global: &config.Global{
			UpdaterName:     "custom",
			Interpreter:     "/bin/bash",
			Script:          "-s",
			ScriptArguments: []string{"<<<", TEST_SCRIPT, "arg2"},
			Timeout:         time.Second * 5,
		},
	}
	customUpdater := updater.New(customConfig)

	err := customUpdater.Update("myservice", "thisversion")

	assert.NoError(t, err)
}

func TestTimeout(t *testing.T) {
	customConfig := &config.Config{
		Global: &config.Global{
			UpdaterName:     "custom",
			Interpreter:     "top",
			Script:          "",
			ScriptArguments: []string{},
			Timeout:         time.Second * 1,
		},
	}
	customUpdater := updater.New(customConfig)

	err := customUpdater.Update("myservice", "thisversion")

	assert.Error(t, err)
}
