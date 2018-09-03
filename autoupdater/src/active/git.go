package active

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/src/logger"
	"github.com/skycoin/services/autoupdater/src/updater"
	"github.com/skycoin/services/autoupdater/store/services"
)

type git struct {
	// Url should be in the format /:owner/:Repository
	url       string
	service   string
	interval  time.Duration
	ticker    *time.Ticker
	lock      sync.Mutex
	tag       string
	date      *time.Time
	updater   updater.Updater
	exit      chan int
	retryTime time.Duration
	retries   int
	log *logger.Logger
	config.CustomLock
}

func newGit(u updater.Updater, service, url string, retries int, retryTime time.Duration, log *logger.Logger) *git {
	retrievedStatus := services.GetStore().Get(service)
	log.Infof("retrieved status %+v", retrievedStatus)
	date := retrievedStatus.LastUpdated.Time

	return &git{
		url:       "https://api.github.com/repos" + url,
		tag:       "0.0.0",
		date:      &date,
		exit:      make(chan int),
		updater:   u,
		service:   service,
		retries:   retries,
		retryTime: retryTime,
		log: log,
	}
}

func (g *git) SetLastRelease(tag string, date *time.Time) {
	g.tag = tag

	if date != nil {
		g.date = date
	}
}

func (g *git) SetInterval(t time.Duration) {
	g.interval = t

	g.lock.Lock()
	if g.ticker != nil {
		g.ticker = time.NewTicker(g.interval)
	}
	g.lock.Unlock()
}

func (g *git) Start() {
	g.ticker = time.NewTicker(g.interval)
	go func() {
		for {
			select {
			case t := <-g.ticker.C:
				g.log.Info("looking for new version at: ", t)
				// Try to fetch new version
				go g.checkIfNew()
			}
		}
	}()
	<-g.exit
}

func (g *git) Stop() {
	g.ticker.Stop()
	g.exit <- 1
}

type ReleaseJSON struct {
	Url string `json:"Url"`
	//Name encodes the name of the release, or its version
	Name        string `json:"Name"`
	PublishedAt string `json:"published_at"`
}

func (g *git) checkIfNew() {
	if g.IsLock() {
		g.log.Warnf("service %s is already being updated... waiting for it to finish", g.service)
	}
	g.Lock()
	defer g.Unlock()

	release := g.fetchGithubRelease()
	publishedTime := parsePublishedTime(release, g.log)

	if g.date.Before(publishedTime) {
		g.log.Info("new version: ", release.Url, ". Published at: ", release.PublishedAt)
		err := g.tryUpdate(release)

		if err != nil {
			logrus.Error(err)
		} else {
			g.storeLastUpdated(publishedTime)
		}
	} else {
		g.log.Info("no new version")
	}
}

func (g *git) fetchGithubRelease() *ReleaseJSON {
	resp, err := http.Get(g.url + "/releases/latest")
	if err != nil {
		g.log.Fatal("cannot contact api, err ", err)
	}
	defer resp.Body.Close()
	release := &ReleaseJSON{}
	err = json.NewDecoder(resp.Body).Decode(release)
	if err != nil {
		g.log.Fatal("cannot unmarshal to a release object, err: ", err)
	}
	if release.PublishedAt == "" {
		g.log.Fatalf("unable to retrieve published at information from %s",
			g.url+"/release/latest. Make sure that the configuration repository exists")
	}

	return release
}

func parsePublishedTime(release *ReleaseJSON, log *logger.Logger) time.Time {
	publishedTime, err := time.Parse(time.RFC3339, release.PublishedAt)
	if err != nil {
		log.Fatal("cannot parse git release date: ", release.PublishedAt, " err: ", err)
	}
	return publishedTime
}

func (g *git) tryUpdate(release *ReleaseJSON) error {
	for i := 0; i < g.retries; i++ {
		err := <-g.updater.Update(g.service, release.Name, g.log)
		if err != nil {
			g.log.Errorf("error on update %s", err)

			if i == (g.retries - 1) {
				return fmt.Errorf("maximum retries attempted, service %s failed to update", g.service)
			} else {
				g.log.Infof("retry again in %s", g.retryTime.String())
			}
		} else {
			break
		}

		time.Sleep(g.retryTime)
	}
	return nil
}

func (g *git) storeLastUpdated(publishedTime time.Time) {
	g.date = &publishedTime
	storeService := services.Service{
		Name:        g.service,
		LastUpdated: services.NewTimeJSON(publishedTime),
	}
	services.GetStore().Store(&storeService)
}
