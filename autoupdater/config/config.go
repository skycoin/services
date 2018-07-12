package config

import (
	"github.com/spf13/viper"
	"path/filepath"
	"github.com/sirupsen/logrus"
	"strings"
)

const CONFIG_TYPE = "toml"
const SERVICES_KEY = "services"

type Config struct {
	Services map[string]string
}

func NewConfig(path string) *Config {
	dir, file := filepath.Split(path)
	cleanFile := strings.TrimSuffix(file, filepath.Ext(file))
	viper.SetConfigType(CONFIG_TYPE)
	viper.SetConfigName(cleanFile)
	viper.AddConfigPath(dir)

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("Unable to read configuration file: %s", err)
	}

	c := &Config{}
	c.parseServices()

	return c
}

func (c *Config) parseServices() {
	c.Services = viper.GetStringMapString(SERVICES_KEY)
}
