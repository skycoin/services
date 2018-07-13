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

type Config struct {
	Global   *Global
	Active   *Active
	Passive  *Passive
	Services map[string]string
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

type Active struct {
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
	c := &Config{}

	if path != "" {
		c.loadConfigFromFile(path)
	}

	return c
}

func (c *Config) loadConfigFromFile(path string) {
	c.Global = &Global{}
	c.Active = &Active{}
	c.Passive = &Passive{}

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
	c.Services = viper.GetStringMapString(SERVICES_KEY)
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
