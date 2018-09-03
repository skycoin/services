package active

import (
	"time"

	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/src/logger"
	"github.com/skycoin/services/autoupdater/src/updater"
)

type Fetcher interface {
	SetInterval(duration time.Duration)
	Start()
	Stop()
}

func New(service string, c *config.Configuration, updater updater.Updater, log *logger.Logger) Fetcher {
	serviceConfig := c.Services[service]
	updateCheckerConfig := c.ActiveUpdateCheckers[serviceConfig.ActiveUpdateChecker]
	updaterConfig := c.Updaters[serviceConfig.Updater]

	if updaterConfig.Kind == "swarm" {
		log.Info("Swarm mode cannot fetch from Git, falling back to Dockerhub")
		updateCheckerConfig.Kind = "dockerhub"
	}

	retryTime, err := time.ParseDuration(updateCheckerConfig.RetryTime)
	if err != nil {
		log.Fatalf("cannot parse retry time %s of service configuration %s", updateCheckerConfig.RetryTime, service)
	}

	switch updateCheckerConfig.Kind {
	case "git":
		return newGit(updater, service, serviceConfig.Repository, updateCheckerConfig.Retries,
			retryTime, log)
	case "naive":
		return newNaive(updater, service, serviceConfig.Repository, updateCheckerConfig.Retries,
			retryTime, log)
	case "dockerhub":
		return newDockerHub(updater, serviceConfig.Repository, serviceConfig.CheckTag,
			serviceConfig.OfficialName, "", log)
	}
	return newDockerHub(updater, serviceConfig.Repository, serviceConfig.CheckTag,
		serviceConfig.OfficialName, "",log)
}
