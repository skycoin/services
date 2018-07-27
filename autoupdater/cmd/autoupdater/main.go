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

const defaultTopic = "top"

var updaterNameAux string

func cmd() *cli.App {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "path to toml configuration file",
			Value:  "",
			EnvVar: "CONFIG",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "swarm",
			Usage: "autoupdates swarm",
			Subcommands: []cli.Command{
				activeCommand("swarm"),
				passiveCommand("swarm"),
			},
		},
		{
			Name:  "custom",
			Usage: "on update notification launches an user provided custom script",
			Flags: customFlags(),
			Subcommands: []cli.Command{
				activeCommand("custom"),
				passiveCommand("custom"),
			},
		},
	}
	return app
}

func customFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "interpreter, i",
			Value:  "/bin/bash",
			Usage:  "interpreter for the custom script",
			EnvVar: "INTERPRETER",
		},
		cli.StringFlag{
			Name:   "script, s",
			Value:  "/etc/skycoin/autoupdater/update.sh",
			Usage:  "custom script to launch on update notification",
			EnvVar: "SCRIPT",
		},
		cli.StringSliceFlag{
			Name:   "arguments,a",
			Value:  &cli.StringSlice{},
			Usage:  "arguments for the script",
			EnvVar: "SCRIPT_ARGUMENTS",
		},
		cli.DurationFlag{
			Name:   "timeout,t",
			Value:  config.DefaultTimeout,
			Usage:  "timeout for the custom script, after which to retry",
			EnvVar: "SCRIPT_TIMEOUT",
		},
	}
}

func passiveCommand(updaterName string) cli.Command {
	return cli.Command{
		Name:  "passive",
		Usage: "waits for update notification",
		Before: func(c *cli.Context) error {
			updaterNameAux = updaterName
			return nil
		},
		Action: passiveAction,
		Flags:  passiveFlags(),
	}
}

func activeCommand(updaterName string) cli.Command {
	return cli.Command{
		Name:  "active",
		Usage: "periodically checks if there is a new version",
		Before: func(c *cli.Context) error {
			updaterNameAux = updaterName
			return nil
		},
		Action: activeAction,
		Flags:  activeFlags(),
	}
}

func passiveAction(c *cli.Context) {
	logrus.Info("passive -> updater: ", c.GlobalString("updater"),
		" broker: ", c.String("message-broker"))

	conf := config.NewConfig(c.GlobalString("config"))

	conf.Global.UpdaterName = updaterNameAux
	conf.Global.Interpreter = stringPickNonZero(conf.Global.Interpreter, c.Parent().String("interpreter"))
	conf.Global.Script = stringPickNonZero(conf.Global.Script, c.Parent().String("script"))
	conf.Global.ScriptArguments = stringSlicePickNonZero(conf.Global.ScriptArguments, c.Parent().StringSlice("arguments"))
	conf.Global.Timeout = durationPickNonZero(conf.Global.Timeout, c.Parent().Duration("timeout"))

	conf.Passive.Urls = stringSlicePickNonZero(conf.Passive.Urls, c.StringSlice("urls"))
	conf.Passive.MessageBroker = stringPickNonZero(conf.Passive.MessageBroker, c.String("message-broker"))

	updaterInstance := updater.New(conf)
	sub := subscriber.New(conf, updaterInstance)
	sub.Subscribe(defaultTopic)
}

func passiveFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "message-broker, m",
			Value:  "nats",
			Usage:  "supported brokers: nats",
			EnvVar: "MESSAGE_BROKER",
		},
		cli.StringSliceFlag{
			Name:   "urls, u",
			Value:  &cli.StringSlice{"http://localhost:2222"},
			Usage:  "urls for the message broker",
			EnvVar: "PASSIVE_URLS",
		},
	}
}

func activeAction(c *cli.Context) {

	conf := config.NewConfig(c.GlobalString("config"))

	conf.Global.UpdaterName = updaterNameAux

	conf.Services[c.String("service")] = &config.Service{
		LocalName:            c.String("service"),
		OfficialName:         c.String("service"),
		ScriptInterpreter:    c.Parent().String("script"),
		ScriptTimeout:        c.Parent().Duration("timeout"),
		ScriptExtraArguments: c.Parent().StringSlice("arguments"),
	}

	conf.Active.Interval = durationPickNonZero(conf.Active.Interval, c.Duration("interval"))
	conf.Active.Tag = stringPickNonZero(conf.Active.Tag, c.String("version"))
	conf.Active.Repository = stringPickNonZero(conf.Active.Repository, c.String("repository"))
	conf.Active.Name = stringPickNonZero(conf.Active.Name, c.String("fetcher"))
	conf.Active.Service = stringPickNonZero(conf.Active.Service, c.String("service"))
	conf.Active.Retries = intPickNonZero(conf.Active.Retries, c.Int("retries"))
	conf.Active.RetryTime = durationPickNonZero(conf.Active.RetryTime, c.Duration("retry-time"))

	updaterInstance := updater.New(conf)
	fetcher := active.New(conf,updaterInstance)
	logrus.Info("active -> updater: ", c.GlobalString("updater"),
		" interval: ", conf.Active.Interval)
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
		cli.IntFlag{
			Name:   "retries",
			Value:  3,
			Usage:  "number of retries to fetch updates",
			EnvVar: "RETRIES",
		},
		cli.DurationFlag{
			Name: "retry-time",
			Value: time.Minute * 5,
			Usage: "time between retries",
			EnvVar: "RETRY_TIME",
		},
	}
}
func main() {
	err := cmd().Run(os.Args)
	if err != nil {
		logrus.Fatal("error running cmd", err)
	}
}

func stringPickNonZero(confValue, flagValue string) string {
	if confValue == "" {
		return flagValue
	}

	return confValue
}

func stringSlicePickNonZero(confValue, flagValue []string) []string {
	if confValue == nil {
		return flagValue
	}

	return confValue
}

func intPickNonZero(confValue, flagValue int) int{
	if confValue == 0 && flagValue != 0 {
		return flagValue
	}

	return confValue
}

func durationPickNonZero(confValue, flagValue time.Duration) time.Duration {
	if confValue == 0 {
		return flagValue
	}

	return confValue
}
