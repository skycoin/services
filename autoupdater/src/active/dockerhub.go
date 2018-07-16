package active

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	"github.com/skycoin/services/autoupdater/src/updater"
)

const SCHEMA_VERSION_HEADER = "application/vnd.docker.distribution.manifest.v2+json"
const URI = "/manifests/latest"
const TOKEN_TEMPLATE = "https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull"

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

func NewDockerHub(updater updater.Updater, repository, tag, service, currentDigest string) *Dockerhub {
	if currentDigest == "" {
		imageName := repository + ":" + tag
		currentDigest = getCurrentDockerImageDigest(imageName)
	}
	parsedRepo := strings.Replace(repository, "/", "", 1)
	logrus.Infof("Retrieved ID is: %s", currentDigest)

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
		TokenTemplate: TOKEN_TEMPLATE,
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
	logrus.Info("Looking for new version at: ", t)
	// Try to fetch new version
	err := g.checkIfNew()
	if err != nil {
		logrus.Info("Cannot contact Dockerhub api: ", err)
		if time.Now().After(g.token.ExpirationDate) {
			logrus.Info("Token expired. Requesting new token...")
			g.getToken()
		}
	}
}

// We need to get a token with pull access to the Repository
func (g *Dockerhub) getToken() {
	tokenRequest := fmt.Sprintf(g.TokenTemplate, g.Repository)
	logrus.Infof("Requesting token to %s", tokenRequest)

	resp, err := http.Get(tokenRequest)
	if err != nil {
		logrus.Fatal("Cannot request a token to: ", tokenRequest, " err: ", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(g.token)
	if err != nil {
		logrus.Fatal("Cannot parse token, err: ", err)
	}
	logrus.Info(fmt.Sprintf("Got token %s", g.token.Token))

	date, err := time.Parse(time.RFC3339, g.token.IssuedAt)
	if err != nil {
		logrus.Fatal("Cannot parse token date: ", err)
	}

	expiresIn, err := time.ParseDuration(fmt.Sprintf("%ds", g.token.ExpiresIn))
	if err != nil {
		logrus.Fatal("Cannot parse expires in: ", err)
	}

	g.token.ExpirationDate = date.Add(expiresIn)
}

// We are looking that the latest image digest is different from the local one
func (g *Dockerhub) checkIfNew() error {
	url := g.Url + URI
	logrus.Info("Looking for new versions at ", url)

	release, err := getLatestDigest(g.Client, g.token.Token, url)
	if err != nil {
		return err
	}
	if release.Config.Digest == "" {
		return fmt.Errorf("Response doesn't contains a digest")
	}
	if g.localDigest != release.Config.Digest {
		logrus.Info("New version: ", release.Config.Digest)
		err := g.Updater.Update(g.Service, g.Repository+":"+g.Tag)
		if err != nil {
			return err
		}
		g.localDigest = release.Config.Digest
	}

	return nil
}

func getCurrentDockerImageDigest(imageName string) string {
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		logrus.Fatal("Cannot connect to docker daemon: ", err)
	}
	image, err := client.InspectImage(imageName)
	if err != nil {
		logrus.Fatal("Cannot inspect local image ", imageName, " err: ", err)
	}
	return image.ID
}

func getLatestDigest(client *http.Client, token, url string) (*DockerReleaseJSON, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Fatal("Cannot create request object, err: ", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", SCHEMA_VERSION_HEADER)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	release := &DockerReleaseJSON{}
	err = json.NewDecoder(resp.Body).Decode(release)
	if err != nil {
		logrus.Fatal("Cannot unmarshal to a release object, err: ", err)
	}

	return release, nil
}
