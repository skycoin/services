package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const jsonFile = `{
  "skycoin": {
	"name": "skycoin",
    "last_updated": "2018-04-26T02:06:00Z"
  },
  "skywire": {
	"name": "skywire",
    "last_updated": "2018-06-26T02:06:00Z"
  }
}`

func TestDecodeFile(t *testing.T) {
	services := decodeServices([]byte(jsonFile))

	assert.NotNil(t, services)
}

func TestGet(t *testing.T) {
	const serviceName = "skycoin"
	serviceLastUpdated := "2018-04-26T02:06:00Z"
	jsonStore := jsonStore{
		services: decodeServices([]byte(jsonFile)),
	}

	service := jsonStore.Get(serviceName)

	assert.Equal(t,service.LastUpdated.Format(time.RFC3339), serviceLastUpdated)
}