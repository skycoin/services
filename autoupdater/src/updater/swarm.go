package updater

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	"fmt"
)

type swarmUpdater struct {
	client *docker.Client
}

func newSwarmUpdater() *swarmUpdater{
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		logrus.Fatal("Cannot connect to the docker daemon: ", err)
	}

	return &swarmUpdater{client}
}

func (s *swarmUpdater) Update(service string, version string) error {
	serviceInfo, err := s.client.InspectService(service)
	if err != nil{
		return fmt.Errorf("Failed to inspect service %s. %s",service, err)
	}

	if serviceInfo.Spec.TaskTemplate.ContainerSpec.Image != version {
		serviceInfo.Spec.TaskTemplate.ContainerSpec.Image = version

		logrus.Info("old index ", serviceInfo.Version.Index)
		updateOptions := docker.UpdateServiceOptions{
			Version:     serviceInfo.Version.Index,
			ServiceSpec: serviceInfo.Spec,
		}
		err = s.client.UpdateService(service, updateOptions)
		if err != nil {
			return fmt.Errorf("Unable to update service %s to version %s. %s",service,version,err )
		}
	}
	return nil
}