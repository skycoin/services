package updater

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/store/services"
)

type Updater interface {
	Update(service, version string) chan error
}

func New(conf *config.Config) Updater {
	services.InitStorer("json")

	normalized := strings.ToLower(conf.Global.UpdaterName)
	logrus.Infof("updater: %s", normalized)

	switch normalized {
	case "swarm":
		return newSwarmUpdater(conf)
	case "custom":
		return newCustomUpdater(conf)
	}

	return newSwarmUpdater(conf)
}
