package updater

import (
	"fmt"

	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/src/logger"
)

type swarmUpdater struct {
	client   *docker.Client
	services map[string]config.ServiceConfig
}

func newSwarmUpdater(services map[string]config.ServiceConfig) *swarmUpdater {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		logrus.Fatal("Cannot connect to the docker daemon: ", err)
	}

	return &swarmUpdater{client, services}
}

func (s *swarmUpdater) Update(service string, version string, log *logger.Logger) chan error {
	errCh := make(chan error)
	go s.update(service,version,log,errCh)
	return errCh
}

func (s *swarmUpdater) update(service string, version string, log *logger.Logger, errCh chan error) {
	localService := s.services[service].LocalName
	serviceInfo, err := s.client.InspectService(localService)
	if err != nil {
		errCh <-fmt.Errorf("failed to inspect service %s. %s", localService, err)
	}

	if serviceInfo.Spec.TaskTemplate.ContainerSpec.Image != version {
		serviceInfo.Spec.TaskTemplate.ContainerSpec.Image = version

		updateOptions := docker.UpdateServiceOptions{
			Version:     serviceInfo.Version.Index,
			ServiceSpec: serviceInfo.Spec,
		}
		log.Info("updating...")
		err = s.client.UpdateService(localService, updateOptions)
		if err != nil {
			errCh <- fmt.Errorf("unable to update service %s to version %s. %s", localService, version, err)
		}
	}

	errCh <- nil
}
