package active

import (
	"time"

	"github.com/skycoin/services/autoupdater/config"
)

type Fetcher interface {
	SetInterval(duration time.Duration)
	Start()
	Stop()
}

func New(c *config.Config) Fetcher {
	switch c.Active.Name {
	case "git":
		return newGit(c.Active.Repository)
	case "dockerhub":
		return NewDockerHub(c.Global.Updater, c.Active.Repository, c.Active.Tag, c.Active.Service, c.Active.CurrentVersion)
	}
	return NewDockerHub(c.Global.Updater, c.Active.Repository, c.Active.Tag, c.Active.Service, c.Active.CurrentVersion)
}
