package subscriber

import (
	"strings"

	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/src/updater"
)

type Subscriber interface {
	Subscribe(topic string)
	Stop()
}

func New(config *config.Config) Subscriber {
	config.Passive.MessageBroker = strings.ToLower(config.Passive.MessageBroker)
	updater := updater.New(config)
	switch config.Passive.MessageBroker{
	case "nats":
		return newNats(updater, config.Passive.Urls[0])
	}

	return newNats(updater, config.Passive.Urls[0])
}
