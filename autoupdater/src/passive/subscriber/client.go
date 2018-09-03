package subscriber

import (
	"strings"

	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/src/logger"
	"github.com/skycoin/services/autoupdater/src/updater"
)

type Subscriber interface {
	Subscribe(topic string)
	Start()
	Stop()
}

func New(config config.SubscriberConfig, updater updater.Updater, log *logger.Logger) Subscriber {
	config.MessageBroker = strings.ToLower(config.MessageBroker)
	switch config.MessageBroker {
	case "nats":
		return newNats(updater, config.Urls[0], log)
	}

	return newNats(updater, config.Urls[0],log)
}
