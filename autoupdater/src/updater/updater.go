package updater

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/src/logger"
)

type Updater interface {
	Update(service, version string, log *logger.Logger) chan error
}

func New(kind string, conf config.Configuration) Updater {

	normalized := strings.ToLower(kind)
	logrus.Infof("updater: %s", normalized)

	switch normalized {
	case "swarm":
		return newSwarmUpdater(conf.Services)
	case "custom":
		return newCustomUpdater(conf.Services)
	}

	return newSwarmUpdater(conf.Services)
}
