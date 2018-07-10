package main

import (
	"github.com/urfave/cli"
	"os"
	"github.com/sirupsen/logrus"
	"time"
	"github.com/skycoin/services/autoupdater/src/passive/subscriber"
	"github.com/skycoin/services/autoupdater/src/updater"
	"github.com/skycoin/services/autoupdater/src/active"
)

type config struct {
	mode string

}

const DEFAULT_URL = "http://localhost:4222"
const DEFAULT_GIT = "/skycoin/skycoin"
const DEFAULT_TOPIC = "top"

func cmd() *cli.App{
	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "updater, u",
			Usage: "defines how to update the software",
			Value: "swarm",
			EnvVar: "UPDATER",
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "passive",
			Usage: "waits for update notification",
			Action: passiveAction,
			Flags: passiveFlags(),
		},
		{
			Name: "active",
			Usage: "periodically checks if there is a new version",
			Action: activeAction,
			Flags: activeFlags(),
		},
	}
	return app
	}

func passiveAction(c *cli.Context){
	logrus.Info("passive -> updater: ",c.GlobalString("updater"),
		" broker: ",c.String("message-broker"))

	subConfig := &subscriber.Config{
		Urls: []string{DEFAULT_URL},
		Updater: updater.New(c.String("updater")),
		Name: c.String("message-broker"),
	}
	sub := subscriber.New(subConfig)
	sub.Subscribe(DEFAULT_TOPIC)
}

func passiveFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name: "message-broker, m",
			Value: "nats",
			Usage: "supported brokers: nats",
			EnvVar: "MESSAGE_BROKER",
		},
	}
}

func activeAction(c *cli.Context){
	logrus.Info("active -> updater: ",c.GlobalString("updater"),
		" interval: ",c.Duration("interval").String())
	fetcher := active.New(c.String("updater"),DEFAULT_GIT)
	fetcher.SetInterval(c.Duration("interval"))
	fetcher.Start()
}

func activeFlags() []cli.Flag {
	return []cli.Flag{
		cli.DurationFlag{
			Name: "interval, i",
			Value: 1*time.Hour,
			Usage: "time interval to check for new version",
			EnvVar: "INTERVAL",
		},
	}
}
func main() {
	err := cmd().Run(os.Args)
	if err != nil{
		logrus.Fatal("error running cmd",err)
	}
}