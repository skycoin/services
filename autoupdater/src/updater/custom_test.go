package updater_test

import (
	"testing"
	"github.com/skycoin/services/autoupdater/src/updater"
	"github.com/skycoin/services/autoupdater/config"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestCustom(t *testing.T) {
	customConfig := &config.Config{
		Global: &config.Global{
			UpdaterName: "custom",
			Interpreter: "/bin/bash",
			Script: "custom_example/custom_script.sh",
			ScriptArguments: []string{"arg1","arg2"},
			Timeout: time.Second * 5,
		},
	}
	customUpdater := updater.New(customConfig)

	err := customUpdater.Update("myservice","thisversion")

	assert.NoError(t,err)
}

func TestTimeout(t *testing.T) {
	customConfig := &config.Config{
		Global: &config.Global{
			UpdaterName: "custom",
			Interpreter: "top",
			Script: "",
			ScriptArguments: []string{},
			Timeout: time.Second * 1,
		},
	}
	customUpdater := updater.New(customConfig)

	err := customUpdater.Update("myservice","thisversion")

	assert.Error(t, err)
}
