package config_test

import (
	"testing"
	"time"

	"github.com/skycoin/services/autoupdater/config"
	"github.com/stretchr/testify/assert"
)

func TestServices(t *testing.T) {
	var EXPECTED_SERVICE_MAP = map[string]config.Service{
		"skycoin/skycoin": config.Service{
			LocalName:            "library/mariadb",
			OfficialName:         "skycoin/skycoin",
			UpdateScript:         "./updater/custom_example/custom_script.sh",
			ScriptInterpreter:    "/bin/bash",
			ScriptTimeoutString:  "5s",
			ScriptTimeout:        5 * time.Second,
			ScriptExtraArguments: []string{"-a 1"},
		},
		"top": config.Service{
			LocalName:            "skywire",
			OfficialName:         "top",
			UpdateScript:         "./updater/custom_example/custom_script.sh",
			ScriptInterpreter:    "/bin/bash",
			ScriptTimeoutString:  "5s",
			ScriptTimeout:        5 * time.Second,
			ScriptExtraArguments: []string{""},
		},
		"sky-node": config.Service{
			LocalName:            "skycoin",
			OfficialName:         "sky-node",
			UpdateScript:         "./updater/custom_example/custom_script.sh",
			ScriptInterpreter:    "/bin/bash",
			ScriptTimeoutString:  "7s",
			ScriptTimeout:        7 * time.Second,
			ScriptExtraArguments: []string{"-a 1", "-b 2"},
		},
		"skywire": config.Service{
			LocalName:            "mystack_skywire",
			OfficialName:         "skywire",
			UpdateScript:         "./updater/custom_example/custom_script.sh",
			ScriptInterpreter:    "/bin/bash",
			ScriptTimeoutString:  "10s",
			ScriptTimeout:        10 * time.Second,
			ScriptExtraArguments: []string{""},
		},
	}

	c := config.NewConfig("../configuration.toml")

	assert.Equal(t, EXPECTED_SERVICE_MAP, c.Services)
}

func TestGlobal(t *testing.T) {
	EXPECTED_GLOBAL := &config.Global{
		UpdaterName: "swarm",
	}

	c := config.NewConfig("../configuration.example.toml")

	assert.Equal(t, EXPECTED_GLOBAL, c.Global)
}

func TestActive(t *testing.T) {
	EXPECTED_ACTIVE := &config.Active{
		Interval:   time.Duration(time.Hour),
		Tag:        "latest",
		Repository: "library/mariadb",
		Name:       "dockerhub",
		Service:    "maria",
	}

	c := config.NewConfig("../configuration.example.toml")

	assert.Equal(t, EXPECTED_ACTIVE, c.Active)
}

func TestPassive(t *testing.T) {
	EXPECTED_PASSIVE := &config.Passive{
		MessageBroker: "nats",
		Urls:          []string{"url1", "url2"},
	}

	c := config.NewConfig("../configuration.example.toml")

	assert.Equal(t, EXPECTED_PASSIVE, c.Passive)
}
