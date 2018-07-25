package updater

import (
	"fmt"

	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	"github.com/skycoin/services/autoupdater/config"
)

type swarmUpdater struct {
	client   *docker.Client
	services map[string]*config.Service
}

func newSwarmUpdater(c *config.Config) *swarmUpdater {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		logrus.Fatal("Cannot connect to the docker daemon: ", err)
	}

	return &swarmUpdater{client, c.Services}
}

func (s *swarmUpdater) Update(service string, version string) error {
	localService := s.services[service].LocalName
	serviceInfo, err := s.client.InspectService(localService)
	if err != nil {
		return fmt.Errorf("failed to inspect service %s. %s", localService, err)
	}

	if serviceInfo.Spec.TaskTemplate.ContainerSpec.Image != version {
		serviceInfo.Spec.TaskTemplate.ContainerSpec.Image = version

		logrus.Info("old index ", serviceInfo.Version.Index)
		updateOptions := docker.UpdateServiceOptions{
			Version:     serviceInfo.Version.Index,
			ServiceSpec: serviceInfo.Spec,
		}
		err = s.client.UpdateService(localService, updateOptions)
		if err != nil {
			return fmt.Errorf("unable to update service %s to version %s. %s", localService, version, err)
		}
	}
	return nil
}
