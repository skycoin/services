package config

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const CONFIG_TYPE = "toml"
const SERVICES_KEY = "services"

var DEFAULT_TIMEOUT time.Duration = 10 * time.Second

const DEFAULT_TIMEOUT_STRING = "10m"

type Config struct {
	Global   *Global
	Active   *Active
	Passive  *Passive
	Services map[string]*Service
}

type Global struct {
	UpdaterName     string
	Interpreter     string
	Script          string
	ScriptArguments []string
	Timeout         time.Duration
}

type Passive struct {
	MessageBroker string
	Urls          []string
}

type Service struct {
	OfficialName         string   `mapstructure:"official_name"`
	LocalName            string   `mapstructure:"local_name"`
	UpdateScript         string   `mapstructure:"update_script"`
	ScriptTimeoutString  string   `mapstructure:"script_timeout"`
	ScriptInterpreter    string   `mapstructure:"script_interpreter"`
	ScriptExtraArguments []string `mapstructure:"script_extra_arguments"`
	ScriptTimeout        time.Duration
	CustomLock
}

type Active struct {
	// Interval in which to look for updates
	Interval time.Duration
	// Fetcher name: Dockerhub or git
	Name string
	// Repository name in the format /:owner/:image, without Tag
	Repository string
	// Image Tag in which to look for updates
	Tag string
	// Service name to update
	Service string
	// Current version of the service
	CurrentVersion string
}

func NewConfig(path string) *Config {
	c := &Config{
		Global:  &Global{},
		Active:  &Active{},
		Passive: &Passive{},
	}

	if path != "" {
		c.loadConfigFromFile(path)
	}

	return c
}

func (c *Config) loadConfigFromFile(path string) {
	dir, file := filepath.Split(path)
	cleanFile := strings.TrimSuffix(file, filepath.Ext(file))
	viper.SetConfigType(CONFIG_TYPE)
	viper.SetConfigName(cleanFile)
	viper.AddConfigPath(dir)

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("Unable to read configuration file: %s", err)
	}

	c.parseGlobal()
	c.parseActive()
	c.parsePassive()
	c.parseServices()
}

func (c *Config) parseGlobal() {
	c.Global.UpdaterName = viper.GetString("global.updater")
}

func (c *Config) parseServices() {
	var services []*Service
	c.Services = make(map[string]*Service)

	err := viper.UnmarshalKey("service", &services)
	if err != nil {
		logrus.Fatalf("Cannot unmarshal toml services configuration %s", err)
	}

	for _, service := range services {
		if service.ScriptTimeoutString == "" {
			service.ScriptTimeoutString = DEFAULT_TIMEOUT_STRING
		}
		service.ScriptTimeout, err = time.ParseDuration(service.ScriptTimeoutString)
		if err != nil {
			logrus.Fatalf("Unable to parse script timeout %s", err)
		}

		c.Services[service.OfficialName] = service
	}
}

func (c *Config) parseActive() {
	c.Active.Interval = viper.GetDuration("active.interval")
	c.Active.Name = viper.GetString("active.name")
	c.Active.Repository = viper.GetString("active.repository")
	c.Active.Service = viper.GetString("active.service")
	c.Active.Tag = viper.GetString("active.tag")
}

func (c *Config) parsePassive() {
	c.Passive.MessageBroker = viper.GetString("passive.message-broker")
	c.Passive.Urls = viper.GetStringSlice("passive.urls")
}
