package updater

import (
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type swarmUpdater struct {
	dclient *client.Client
	version swarm.Version
}

func newSwarmUpdater() *swarmUpdater{
	c,err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logrus.Fatal("Unable to connect to docker daemon: ", err)
	}

	swarm, err := c.SwarmInspect(context.Background())
	if err != nil {
		logrus.Fatal("Unable to connect to swarm: ", err)
	}

	return &swarmUpdater{c, swarm.Version}
}

func (s *swarmUpdater) Update(service string) {
	// TODO we need to set a default service name in the client side.
	// The service name from our side can map to a different name defined by the user
	// so it will match their own service name in the swarm
	serviceInfo, _, err := s.dclient.ServiceInspectWithRaw(context.Background(),service,types.ServiceInspectOptions{})
	if err != nil {
		logrus.Fatal("Unable to connect to swarm: ", err)
	}

	serviceUpdateOptions := types.ServiceUpdateOptions{
		//TODO look what to tweak
	}

	resp, err := s.dclient.ServiceUpdate(context.Background(),
			serviceInfo.ID,
			s.version,
			serviceInfo.Spec,
			serviceUpdateOptions,
	)

	// TODO maybe we want to retry?
	if err != nil {
		logrus.Fatal("Unable to connect to update service ",service,": ", err)
	}

	if len(resp.Warnings) > 0 {
		logrus.Warn("Warnings generated on update of service ", service)
		for _, warn := range resp.Warnings {
			logrus.Warn(warn)
		}
	}
}