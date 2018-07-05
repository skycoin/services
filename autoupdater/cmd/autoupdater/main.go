package main

import (
	"github.com/urfave/cli"
	"os"
	"github.com/sirupsen/logrus"
	"time"
)

type config struct {
	mode string

}

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