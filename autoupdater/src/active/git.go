package active

import (
	"time"
	"github.com/sirupsen/logrus"
	"sync"
	"net/http"
	"encoding/json"
)

type git struct {
	// Url should be in the format /:owner/:Repository
	url      string
	interval time.Duration
	ticker   *time.Ticker
	lock sync.Mutex
	tag string
	date *time.Time
	exit chan int
}

func newGit(url string) *git {
	date := time.Date(1999,10,1,1,1,1,1, time.UTC)
	return &git{url: "https://api.github.com/repos"+url,tag: "0.0.0",date: &date,exit: make(chan int),}
}

func (g *git) SetLastRelease(tag string,  date *time.Time) {
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
				logrus.Info("Looking for new version at: ", t)
				// Try to fetch new version
				g.checkIfNew()
			}
		}
	}()
	<- g.exit
}

func (g *git) Stop() {
	g.ticker.Stop()
	g.exit <- 1
}

type ReleaseJSON struct {
	Url string `json:"Url"`
	PublishedAt string `json:"published_at"`
}

func (g *git) checkIfNew() {
	resp, err := http.Get(g.url+"/releases/latest")
	if err != nil {
		logrus.Fatal("cannot contact api ", g.url+"/releases/latest", " err ", err)
	}
	defer resp.Body.Close()

	release := &ReleaseJSON{}
	err = json.NewDecoder(resp.Body).Decode(release)
	if err != nil {
		logrus.Fatal("cannot unmarshal to a release object, err: ", err)
	}

	publishedTime, err := time.Parse(time.RFC3339, release.PublishedAt)
	if err != nil {
		logrus.Fatal("Cannot parse git release date: ", release.PublishedAt, " err: ", err)
	}

	if g.date.Before(publishedTime) {
		logrus.Info("New version: ",release.Url,". Published at: ", release.PublishedAt)
	}
}
