package updater

import (
	"strings"

	"github.com/skycoin/services/autoupdater/config"
)

type Updater interface {
	Update(service, version string) error
}

func New(conf *config.Config) Updater {
	normalized := strings.ToLower(conf.Global.UpdaterName)

	switch normalized {
	case "swarm":
		return newSwarmUpdater()
	case "custom":
		return newCustomUpdater(conf.Global)
	}

	return newSwarmUpdater()
}
