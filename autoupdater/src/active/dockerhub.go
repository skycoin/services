package active

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
)

const SCHEMA_VERSION_HEADER = "application/vnd.docker.distribution.manifest.v2+json"
const URI = "/manifests/latest"

type dockerhub struct {
	// url should be in the format /:owner/:repository
	url         string
	repository  string
	client      *http.Client
	interval    time.Duration
	ticker      *time.Ticker
	lock        sync.Mutex
	tag         string
	localDigest string
	exit        chan int
	token       *DockerHubToken
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

func newDockerHub(url string) *dockerhub {
	latest := "latest"
	imageName := url + ":" + latest
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		logrus.Fatal("Cannot connect to docker daemon: ", err)
	}

	image, err := client.InspectImage(imageName)
	if err != nil {
		logrus.Fatal("Cannot inspect local image ", imageName, " err: ", err)
	}

	return &dockerhub{
		url:         "https://registry.hub.docker.com/v2" + url,
		repository:  strings.Replace(url, "/", "", 1),
		client:      &http.Client{},
		tag:         latest,
		localDigest: image.ID,
		exit:        make(chan int),
		token:       &DockerHubToken{},
	}
}

func (g *dockerhub) SetLastRelease(tag string, date *time.Time) {
	g.tag = tag
}

func (g *dockerhub) SetInterval(t time.Duration) {
	g.interval = t

	g.lock.Lock()
	if g.ticker != nil {
		g.ticker = time.NewTicker(g.interval)
	}
	g.lock.Unlock()
}

func (g *dockerhub) Start() {
	g.ticker = time.NewTicker(g.interval)
	g.getToken()
	go func() {
		for {
			select {
			case t := <-g.ticker.C:
				logrus.Info("Looking for new version at: ", t)
				// Try to fetch new version
				err := g.checkIfNew()
				if err != nil {
					logrus.Info("Cannot contact dockerhub api: ", err)
					if time.Now().After(g.token.ExpirationDate) {
						logrus.Info("Token expired. Requesting new token...")
						g.getToken()
					}
				}
			}
		}
	}()
	<-g.exit
}

func (g *dockerhub) Stop() {
	g.ticker.Stop()
	g.exit <- 1
}

// We need to get a token with pull access to the repository
// How do we make sure our token is still valid?
// 1) On request to the API we can check if we were unauthorized. In that case
// we can request a new token
func (g *dockerhub) getToken() {
	tokenRequest := fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull", g.repository)
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
func (g *dockerhub) checkIfNew() error {
	url := g.url + URI
	logrus.Info("Performing request to ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Fatal("Cannot create request object, err: ", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", g.token.Token))
	req.Header.Add("Accept", SCHEMA_VERSION_HEADER)
	resp, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	release := &DockerReleaseJSON{}
	err = json.NewDecoder(resp.Body).Decode(release)
	if err != nil {
		logrus.Fatal("Cannot unmarshal to a release object, err: ", err)
	}
	logrus.Infof("Unmarshaled %+v", release)
	if release.Config.Digest == "" {
		return fmt.Errorf("Response doesn't contains a digest")
	}
	if g.localDigest != release.Config.Digest {
		logrus.Info("New version: ", release.Config.Digest)
	}

	return nil
}
