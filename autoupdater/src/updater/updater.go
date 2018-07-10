package updater

import "strings"

type Updater interface {
	Update(service , version string)
}

func New(name string) Updater {
	normalized := strings.ToLower(name)

	switch normalized{
	case "swarm":
		return newSwarmUpdater()
	}

	return newSwarmUpdater()
}