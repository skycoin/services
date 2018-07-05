package main

import "github.com/urfave/cli"

type config struct {
	mode string

}

func cmd() *cli.App{
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name: "passive",
			Usage: "waits for update notification",
			Action: passiveAction,
			Flags: passiveFlags(),
		},
	}
	return app
	}

func passiveAction(c *cli.Context){}

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

func main() {
}