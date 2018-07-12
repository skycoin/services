package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/src/active"
	"github.com/skycoin/services/autoupdater/src/passive/subscriber"
	"github.com/skycoin/services/autoupdater/src/updater"
	"github.com/urfave/cli"
)

const DEFAULT_URL = "http://localhost:4222"
const DEFAULT_TOPIC = "top"

func cmd() *cli.App {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "updater, u",
			Usage:  "defines how to update the software",
			Value:  "swarm",
			EnvVar: "UPDATER",
		},
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "path to toml configuration file",
			Value:  "",
			EnvVar: "CONFIG",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "passive",
			Usage:  "waits for update notification",
			Action: passiveAction,
			Flags:  passiveFlags(),
		},
		{
			Name:   "active",
			Usage:  "periodically checks if there is a new version",
			Action: activeAction,
			Flags:  activeFlags(),
		},
	}
	return app
}

func passiveAction(c *cli.Context) {
	logrus.Info("passive -> updater: ", c.GlobalString("updater"),
		" broker: ", c.String("message-broker"))

	conf := config.NewConfig(c.GlobalString("config"))
	conf.Global.UpdaterName = c.String("updater")
	conf.Global.Updater = updater.New(conf.Global.UpdaterName)
	conf.Passive = &config.Passive{
		Urls:    []string{DEFAULT_URL},
		MessageBroker:    c.String("message-broker"),
	}
	sub := subscriber.New(conf)
	sub.Subscribe(DEFAULT_TOPIC)
}

func passiveFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "message-broker, m",
			Value:  "nats",
			Usage:  "supported brokers: nats",
			EnvVar: "MESSAGE_BROKER",
		},
	}
}

func activeAction(c *cli.Context) {
	logrus.Info("active -> updater: ", c.GlobalString("updater"),
		" interval: ", c.Duration("interval").String())

	conf := config.NewConfig(c.GlobalString("config"))
	conf.Global.Updater = updater.New(c.String("updater"))
	conf.Active = &config.Active{
		Interval:	c.Duration("interval"),
		Tag:        c.String("version"),
		Repository: c.String("repository"),
		Name:       c.String("fetcher"),
		Service:    c.String("service"),
	}

	fetcher := active.New(conf)
	fetcher.SetInterval(conf.Active.Interval)
	fetcher.Start()
}

func activeFlags() []cli.Flag {
	return []cli.Flag{
		cli.DurationFlag{
			Name:   "interval, i",
			Value:  1 * time.Hour,
			Usage:  "time interval to check for new version",
			EnvVar: "INTERVAL",
		},
		cli.StringFlag{
			Name:   "repository, r",
			Value:  "/skycoin/skycoin",
			Usage:  "repository to fetch updates from",
			EnvVar: "ACTIVE_REPOSITORY",
		},
		cli.StringFlag{
			Name:   "version, v",
			Value:  "latest",
			Usage:  "software version to look for updates",
			EnvVar: "ACTIVE_VERSION",
		},
		cli.StringFlag{
			Name:   "fetcher, f",
			Value:  "dockerhub",
			Usage:  "fetcher used to look for updates: dockerhub or git",
			EnvVar: "ACTIVE_FETCHER",
		},
		cli.StringFlag{
			Name:   "service, s",
			Value:  "skycoin-node",
			Usage:  "service name to be updated",
			EnvVar: "ACTIVE_SERVICE",
		},
	}
}
func main() {
	err := cmd().Run(os.Args)
	if err != nil {
		logrus.Fatal("error running cmd", err)
	}
}
