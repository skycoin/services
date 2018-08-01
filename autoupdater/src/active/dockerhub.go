package active

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/skycoin/services/autoupdater/config"
	"github.com/skycoin/services/autoupdater/src/logger"
	"github.com/skycoin/services/autoupdater/src/updater"
)

const schemaVersionHeader = "application/vnd.docker.distribution.manifest.v2+json"
const uri = "/manifests/latest"
const tokenTemplate = "https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull"

type Dockerhub struct {
	// Url should be in the format /:owner/:Repository
	Url           string
	Repository    string
	Service       string
	Client        *http.Client
	Interval      time.Duration
	Ticker        *time.Ticker
	lock          sync.Mutex
	Tag           string
	localDigest   string
	exit          chan int
	token         *DockerHubToken
	TokenTemplate string
	Updater       updater.Updater
	log *logger.Logger
	config.CustomLock
}

type DockerHubToken struct {
	Token          string    `json:"token,omitempty"`
	AccessToken    string    `json:"access_token,omitempty"`
	ExpiresIn      int       `json:"expires_in,omitempty"`
	IssuedAt       string    `json:"issued_at,omitempty"`
	ExpirationDate time.Time `json:"-"`
}

type DockerReleaseJSON struct {
	SchemaVersion int              `json:"schemaVersion"`
	Config        DockerConfigJSON `json:"config"`
}

type DockerConfigJSON struct {
	Digest string `json:"digest"`
}

func newDockerHub(updater updater.Updater, repository, tag, service, currentDigest string, l *logger.Logger) *Dockerhub {
	if currentDigest == "" {
		imageName := repository + ":" + tag
		currentDigest = getCurrentDockerImageDigest(imageName,l)
	}
	parsedRepo := strings.Replace(repository, "/", "", 1)
	l.Infof("retrieved ID is: %s", currentDigest)

	return &Dockerhub{
		Url:           "https://registry.hub.docker.com/v2" + repository,
		Repository:    parsedRepo,
		Client:        &http.Client{},
		Tag:           tag,
		localDigest:   currentDigest,
		exit:          make(chan int),
		token:         &DockerHubToken{},
		Updater:       updater,
		Service:       service,
		TokenTemplate: tokenTemplate,
		log:l,
	}
}

func (g *Dockerhub) SetLastRelease(tag string, date *time.Time) {
	g.Tag = tag
}

func (g *Dockerhub) SetInterval(t time.Duration) {
	g.Interval = t

	g.lock.Lock()
	if g.Ticker != nil {
		g.Ticker = time.NewTicker(g.Interval)
	}
	g.lock.Unlock()
}

func (g *Dockerhub) Start() {
	g.Ticker = time.NewTicker(g.Interval)
	g.getToken()
	go func() {
		for {
			select {
			case t := <-g.Ticker.C:
				g.checkUpdate(t)
			}
		}
	}()
	<-g.exit
}

func (g *Dockerhub) Stop() {
	g.Ticker.Stop()
	g.exit <- 1
}

func (g *Dockerhub) checkUpdate(t time.Time) {
	if g.IsLock() {
		g.log.Warn("service %s is already being updated... waiting for it to finish")
	}
	g.Lock()
	defer g.Unlock()
	g.log.Info("looking for new version at: ", t)

	// Try to fetch new version
	err := g.updateIfNew()
	if err != nil {
		g.log.Info("cannot contact Dockerhub api: ", err)
		if time.Now().After(g.token.ExpirationDate) {
			g.log.Info("token expired. Requesting new token...")
			g.getToken()
		}
	}
}

// We need to get a token with pull access to the Repository
func (g *Dockerhub) getToken() {
	tokenRequest := fmt.Sprintf(g.TokenTemplate, g.Repository)
	g.log.Infof("requesting token to %s", tokenRequest)

	resp, err := http.Get(tokenRequest)
	if err != nil {
		g.log.Fatal("cannot request a token to: ", tokenRequest, " err: ", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(g.token)
	if err != nil {
		g.log.Fatal("cannot parse token, err: ", err)
	}
	g.log.Info(fmt.Sprintf("Got token %s", g.token.Token))

	date, err := time.Parse(time.RFC3339, g.token.IssuedAt)
	if err != nil {
		g.log.Fatal("cannot parse token date: ", err)
	}

	expiresIn, err := time.ParseDuration(fmt.Sprintf("%ds", g.token.ExpiresIn))
	if err != nil {
		g.log.Fatal("cannot parse expires in: ", err)
	}

	g.token.ExpirationDate = date.Add(expiresIn)
}

// We are looking that the latest image digest is different from the local one
func (g *Dockerhub) updateIfNew() error {
	url := g.Url + uri
	g.log.Info("looking for new versions at ", url)

	release, err := g.getLatestDigest(g.Client, g.token.Token, url)
	if err != nil {
		return err
	}
	if release.Config.Digest == "" {
		return fmt.Errorf("response doesn't contains a digest")
	}
	if g.localDigest != release.Config.Digest {
		g.log.Info("new version: ", release.Config.Digest)
		err := <- g.Updater.Update(g.Service, g.Repository+":"+g.Tag, g.log)
		if err != nil {
			return err
		}
		g.localDigest = release.Config.Digest
	} else {
		g.log.Info("no new version")
	}

	return nil
}

func getCurrentDockerImageDigest(imageName string, log *logger.Logger) string {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		log.Fatal("cannot connect to docker daemon: ", err)
	}
	image, err := client.InspectImage(imageName)
	if err != nil {
		log.Fatal("cannot inspect local image ", imageName, " err: ", err)
	}
	return image.ID
}

func (g *Dockerhub) getLatestDigest(client *http.Client, token, url string) (*DockerReleaseJSON, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		g.log.Fatal("cannot create request object, err: ", err)
	}

	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("accept", schemaVersionHeader)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	release := &DockerReleaseJSON{}
	err = json.NewDecoder(resp.Body).Decode(release)
	if err != nil {
		unmarshalErr := fmt.Errorf("error %s does the repository exists?", err)
		return nil, unmarshalErr
	}

	return release, nil
}
