package updater

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/services/autoupdater/config"
)

type Updater interface {
	Update(service, version string) error
}

func New(conf *config.Config) Updater {
	normalized := strings.ToLower(conf.Global.UpdaterName)
	logrus.Info("updater: %s", normalized)

	switch normalized {
	case "swarm":
		return newSwarmUpdater(conf)
	case "custom":
		return newCustomUpdater(conf)
	}

	return newSwarmUpdater(conf)
}
