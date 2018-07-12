package config_test

import (
	"testing"
	"github.com/skycoin/services/autoupdater/config"
	"github.com/stretchr/testify/assert"
)


func TestConfig(t *testing.T) {
	var EXPECTED_SERVICE_MAP map[string]string = map[string]string{
		"skycoin/skycoin":"library/mariadb",
		"top" : "skywire",
		"sky-node" : "skycoin",
		"skywire" : "mystack_skywire",
	}

	c := config.NewConfig("../service_mapping_example.toml")

	assert.Equal(t, EXPECTED_SERVICE_MAP,c.Services)
}
