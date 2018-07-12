package config_test

import (
	"testing"
	"time"

	"github.com/skycoin/services/autoupdater/config"
	"github.com/stretchr/testify/assert"
)

func TestServices(t *testing.T) {
	var EXPECTED_SERVICE_MAP = map[string]string{
		"skycoin/skycoin": "library/mariadb",
		"top":             "skywire",
		"sky-node":        "skycoin",
		"skywire":         "mystack_skywire",
	}

	c := config.NewConfig("../service_mapping_example.toml")

	assert.Equal(t, EXPECTED_SERVICE_MAP, c.Services)
}

func TestGlobal(t *testing.T) {
	EXPECTED_GLOBAL:= &config.Global{
		UpdaterName: "swarm",
	}


	c := config.NewConfig("../service_mapping_example.toml")

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

	c := config.NewConfig("../service_mapping_example.toml")

	assert.Equal(t, EXPECTED_ACTIVE, c.Active)
}

func TestPassive(t *testing.T) {
	EXPECTED_PASSIVE := &config.Passive{
		MessageBroker: "nats",
		Urls:          []string{"url1", "url2"},
	}

	c := config.NewConfig("../service_mapping_example.toml")

	assert.Equal(t, EXPECTED_PASSIVE, c.Passive)
}

