package active

import (
	"time"
	"github.com/skycoin/services/autoupdater/src/updater"
)

type Config struct {
	// Fetcher name: dockerhub or git
	Name string

	// Repository name in the format /:owner/:image, without tag
	Repository string

	// Image tag in which to look for updates
	Tag string

	// Service updater
	Updater updater.Updater

	// Service name to update
	Service string
}

type Fetcher interface {
	SetInterval(duration time.Duration)
	Start()
	Stop()
}

func New(c *Config) Fetcher {
	switch c.Name{
	case "git":
		return newGit(c.Repository)
	case "dockerhub":
		return newDockerHub(c.Updater, c.Repository, c.Tag, c.Service)
	}
	return newDockerHub(c.Updater, c.Repository, c.Tag, c.Service)
}
