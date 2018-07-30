package subscriber

import (
	"sync"

	gonats "github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
	"github.com/skycoin/services/autoupdater/src/logger"
	"github.com/skycoin/services/autoupdater/src/updater"
)

type nats struct {
	updater    updater.Updater
	url        string
	connection *gonats.Conn
	closer     chan int
	topic      string
	log *logger.Logger
	sync.Mutex
}

func newNats(u updater.Updater, url string, log *logger.Logger) *nats {
	connection, err := gonats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	return &nats{
		updater:    u,
		url:        url,
		connection: connection,
		closer:     make(chan int),
		log: log,
	}
}

func (n *nats) Subscribe(topic string) {
	n.Lock()
	defer n.Unlock()
	n.topic = topic
}

func (n *nats) Start() {
	n.connection.Subscribe(n.topic, n.onUpdate)
	n.log.Infof("subscribed to %s",n.topic)
	<-n.closer
	n.log.Info("stop")
}

func (n *nats) Stop() {
	n.closer <- 1
}

func (n *nats) onUpdate(msg *gonats.Msg) {
	n.log.Info("received update notification")
	err := <- n.updater.Update(msg.Subject, string(msg.Data), n.log)
	if err != nil {
		logrus.Fatal(err)
	}
}
