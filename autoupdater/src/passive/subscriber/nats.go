package subscriber

import (
	"fmt"

	gonats "github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
	"github.com/skycoin/services/autoupdater/src/updater"
)

type nats struct {
	updater    updater.Updater
	url        string
	connection *gonats.Conn
	closer     chan int
}

func newNats(u updater.Updater, url string) *nats {
	connection, err := gonats.Connect(url)
	if err != nil {
		logrus.Fatal("cannot connect to NATS", err)
	}
	return &nats{u, url, connection, make(chan int)}
}

func (n *nats) Subscribe(topic string) {
	n.connection.Subscribe(topic, n.onUpdate)
	<-n.closer
}

func (n *nats) Stop() {
	n.closer <- 1
}

func (n *nats) onUpdate(msg *gonats.Msg) {
	fmt.Println(string(msg.Data))
	err := n.updater.Update(msg.Subject, string(msg.Data))
	if err != nil {
		logrus.Fatal(err)
	}
}
