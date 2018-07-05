package subscriber

import (
	gonats "github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
	"fmt"
)

type nats struct {
	url string
	connection *gonats.Conn
}

func newNats(url string) *nats{
	connection, err := gonats.Connect(url)
	if err != nil {
		logrus.Error("cannot connect to NATS",err)
	}
	return &nats{url, connection}
}

func (n *nats) Subscribe(topic string) {
	n.connection.Subscribe(topic, onUpdate)
}

func onUpdate(msg *gonats.Msg) {
	//TODO use function to update a service
	fmt.Println(string(msg.Data))
}
