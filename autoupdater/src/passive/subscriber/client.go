package subscriber

import (
	"strings"

	"github.com/skycoin/services/autoupdater/config"
)

type Subscriber interface {
	Subscribe(topic string)
	Stop()
}

func New(config *config.Config) Subscriber {
	config.Passive.MessageBroker = strings.ToLower(config.Passive.MessageBroker)
	switch config.Passive.MessageBroker{
	case "nats":
		return newNats(config.Global.Updater, config.Passive.Urls[0])
	}

	return newNats(config.Global.Updater, config.Passive.Urls[0])
}
