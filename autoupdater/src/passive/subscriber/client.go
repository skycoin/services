package subscriber

import (
	"github.com/skycoin/services/autoupdater/src/updater"
	"strings"
)

type Config struct {
	Name string
	Urls []string
	Updater updater.Updater
}

type Subscriber interface {
	Subscribe(topic string)
	Stop()
}

func New(config *Config) Subscriber {
	config.Name = strings.ToLower(config.Name)
	switch config.Name {
	case "nats":
		return newNats(config.Updater, config.Urls[0])
	}

	return newNats(config.Updater, config.Urls[0])
}