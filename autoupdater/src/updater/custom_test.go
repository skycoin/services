package updater_test

import (
	"testing"
	"time"

	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/src/updater"
	"github.com/stretchr/testify/assert"
)

const testScript = `
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
		},
		Services: map[string]*config.Service{
			"myservice": &config.Service{
				LocalName: "myservice",
				OfficialName: "myservice",
				ScriptInterpreter:     "/bin/bash",
				UpdateScript:          "-s",
				ScriptExtraArguments: []string{"<<<", testScript, "arg2"},
				ScriptTimeout:         time.Second * 5,
			},
		},
	}
	customUpdater := updater.New(customConfig)

	err := <- customUpdater.Update("myservice", "thisversion")

	assert.NoError(t, err)
}

func TestTimeout(t *testing.T) {
	customConfig := &config.Config{
		Global: &config.Global{
			UpdaterName:     "custom",
		},
		Services: map[string]*config.Service{
			"myservice": &config.Service{
				LocalName: "myservice",
				OfficialName: "myservice",
				ScriptInterpreter:     "top",
				UpdateScript:          "",
				ScriptExtraArguments: []string{},
				ScriptTimeout:         time.Second * 1,
			},
		},
	}
	customUpdater := updater.New(customConfig)

	err := <- customUpdater.Update("myservice", "thisversion")

	assert.Error(t, err)
}
