package updater_test

import (

"testing"

"github.com/skycoin/services/autoupdater/config"
"github.com/skycoin/services/autoupdater/src/logger"
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
	customConfig := config.Configuration{
		Updaters: map[string]config.UpdaterConfig{
			"test":{
				Kind: "custom",
			},
		},
		Services: map[string]config.ServiceConfig{
			"myservice": {
				LocalName: "myservice",
				OfficialName: "myservice",
				ScriptInterpreter:     "/bin/bash",
				UpdateScript:          "-s",
				ScriptExtraArguments: []string{"<<<", testScript, "arg2"},
				ScriptTimeout:        "5s",
				Updater: "test",
			},
		},
	}
	customUpdater := updater.New("custom",customConfig)

	log := logger.NewLogger("myservice")
	err := <- customUpdater.Update("myservice", "thisversion",log)

	assert.NoError(t, err)
}

func TestTimeout(t *testing.T) {
	customConfig := config.Configuration{
		Updaters: map[string]config.UpdaterConfig{
			"test":{
				Kind: "custom",
			},
		},
		Services: map[string]config.ServiceConfig{
			"myservice": {
				LocalName: "myservice",
				OfficialName: "myservice",
				Updater: "test",
				ScriptInterpreter:     "top",
				UpdateScript:          "",
				ScriptExtraArguments: []string{},
				ScriptTimeout: "1s",
			},
		},
	}
	customUpdater := updater.New("custom",customConfig)

	log := logger.NewLogger("myservice")
	err := <- customUpdater.Update("myservice", "thisversion", log)

	assert.Error(t, err)
}
