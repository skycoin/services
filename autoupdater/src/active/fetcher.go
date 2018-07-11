package active

import (
	"time"
	"github.com/skycoin/services/autoupdater/src/updater"
)

type Config struct {
	// Fetcher name: Dockerhub or git
	Name string

	// Repository name in the format /:owner/:image, without Tag
	Repository string

	// Image Tag in which to look for updates
	Tag string

	// Service updater
	Updater updater.Updater

	// Service name to update
	Service string

	// Current version of the service
	CurrentVersion string
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
		return NewDockerHub(c.Updater, c.Repository, c.Tag, c.Service, c.CurrentVersion)
	}
	return NewDockerHub(c.Updater, c.Repository, c.Tag, c.Service, c.CurrentVersion)
}

