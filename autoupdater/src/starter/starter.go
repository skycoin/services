package starter

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/src/active"
	"github.com/skycoin/services/autoupdater/src/logger"
	"github.com/skycoin/services/autoupdater/src/passive/subscriber"
	"github.com/skycoin/services/autoupdater/src/updater"
	"github.com/skycoin/services/autoupdater/store/services"
)

type Starter struct {
	activeCheckers map[string]active.Fetcher
	passiveCheckers map[string]subscriber.Subscriber
	updaters map[string]updater.Updater
}

func New(conf config.Configuration) *Starter{
	s := &Starter{
		activeCheckers: map[string]active.Fetcher{},
		passiveCheckers: map[string]subscriber.Subscriber{},
		updaters: map[string]updater.Updater{},
	}

	services.InitStorer("json")

	s.createUpdaters(conf)
	s.createCheckers(conf)

	return s
}

func (s *Starter) Start() {
	for _, checker := range s.activeCheckers {
		go checker.Start()
	}

	for _, checker := range s.passiveCheckers {
		go checker.Start()
	}
}

func (s *Starter) Stop() {
	for _, checker := range s.activeCheckers {
		 checker.Stop()
	}

	for _, checker := range s.passiveCheckers {
		 checker.Stop()
	}
}

func (s *Starter) createUpdaters(conf config.Configuration) {
	for name, c := range conf.Updaters{
		u := updater.New(c.Kind, conf)
		s.updaters[name] = u
	}
}

func (s *Starter) createCheckers(conf config.Configuration) {
	for name, c := range conf.Services{
		if c.ActiveUpdateChecker != "" {
			activeConfig, ok := conf.ActiveUpdateCheckers[c.ActiveUpdateChecker]
			if !ok {
				logrus.Warnf("%s checker not defined for service %s, skipping service",
					c.ActiveUpdateChecker, name)
				continue
			}

			interval, err := time.ParseDuration(activeConfig.Interval)
			if err != nil {
				logrus.Fatalf("cannot parse interval %s of active checker configuration %s. %s", activeConfig.Interval,
					c.ActiveUpdateChecker, err)
			}
			log := logger.NewLogger(name)

			checker := active.New(name, &conf, s.updaters[c.Updater],log)
			checker.SetInterval(interval)
			s.activeCheckers[name] = checker
		} else {
			passiveConfig, ok := conf.PassiveUpdateCheckers[c.PassiveUpdateChecker]
			if !ok {
				logrus.Warnf("%s checker not defined for service %s, skipping service",
					c.ActiveUpdateChecker, name)
				continue
			}
			log := logger.NewLogger(name)

			sub :=  subscriber.New(passiveConfig, s.updaters[c.Updater],log)
			s.passiveCheckers[name] = sub
			sub.Subscribe(passiveConfig.Topic)
		}
	}
}
