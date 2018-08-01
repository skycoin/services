package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Updaters              map[string]UpdaterConfig    `yaml:"updaters"`
	ActiveUpdateCheckers  map[string]FetcherConfig    `yaml:"active_update_checkers"`
	PassiveUpdateCheckers map[string]SubscriberConfig `yaml:"passive_update_checkers"`
	Services              map[string]ServiceConfig    `yaml:"services"`
}

type UpdaterConfig struct {
	Kind string `yaml:"kind"`
}

type FetcherConfig struct {
	Interval  string `yaml:"interval"`
	Retries   int    `yaml:"retries"`
	RetryTime string `yaml:"retry_time"`
	Kind      string `yaml:"kind"`
}

type SubscriberConfig struct {
	MessageBroker string   `yaml:"message-broker"`
	Topic         string   `yaml:"topic"`
	Urls          []string `yaml:"urls"`
}

type ServiceConfig struct {
	OfficialName         string   `yaml:"official_name"`
	LocalName            string   `yaml:"local_name"`
	UpdateScript         string   `yaml:"update_script"`
	ScriptTimeout        string   `yaml:"script_timeout"`
	ScriptInterpreter    string   `yaml:"script_interpreter"`
	ScriptExtraArguments []string `yaml:"script_extra_arguments"`
	ActiveUpdateChecker  string   `yaml:"active_update_checker"`
	PassiveUpdateChecker string   `yaml:"passive_update_checker"`
	CheckTag             string   `yaml:"check_tag"`
	Updater              string   `yaml:"updater"`
	Repository           string   `yaml:"repository"`
}

func New(path string) Configuration {
	confPath := defaultPathIfNil(path)

	b, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic(err)
	}

	conf := Configuration{
		ActiveUpdateCheckers:  make(map[string]FetcherConfig),
		PassiveUpdateCheckers: make(map[string]SubscriberConfig),
		Services:              make(map[string]ServiceConfig),
	}

	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		panic(err)
	}

	return conf
}

func defaultPathIfNil(path string) string {
	if path == "" {
		gopath := os.Getenv("$GOPATH")
		confPath := filepath.Join(gopath, "src", "github.com", "skycoin", "services", "autoupdater",
			"configuration.yml")

		return confPath
	}
	return path
}
