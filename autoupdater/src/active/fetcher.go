package active

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/src/updater"
)

type Fetcher interface {
	SetInterval(duration time.Duration)
	Start()
	Stop()
}

func New(c *config.Config) Fetcher {
	updater := updater.New(c)
	if c.Global.UpdaterName == "swarm" {
		logrus.Info("Swarm mode cannot fetch from Git, falling back to Dockerhub")
		c.Active.Name = "dockerhub"
	}

	switch c.Active.Name {
	case "git":
		return newGit(updater, c.Active.Service, c.Active.Repository)
	case "dockerhub":
		return NewDockerHub(updater, c.Active.Repository, c.Active.Tag, c.Active.Service, c.Active.CurrentVersion)
	}
	return NewDockerHub(updater, c.Active.Repository, c.Active.Tag, c.Active.Service, c.Active.CurrentVersion)
}
