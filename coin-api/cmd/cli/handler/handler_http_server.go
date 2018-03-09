package handler

import (
	"context"

	"github.com/BurntSushi/toml"
	"github.com/labstack/echo"
	"github.com/skycoin/services/coin-api/cmd/servd"
	"github.com/urfave/cli"
	"log"
)

// ServerHTTP is a CLI handler of an HTTP server
type ServerHTTP struct {
	server *echo.Echo
}

//NewServerHTTP returns an http server
func NewServerHTTP() *ServerHTTP {
	return &ServerHTTP{}
}

// Start starts the http server
func (s ServerHTTP) Start(c *cli.Context) error {
	cfgFile := c.Args().First()

	var config = &servd.Config{}
	_, err := toml.DecodeFile(cfgFile, config)

	if err != nil {
		log.Fatal(err)
	}

	srv, err := servd.Start(config)
	if err != nil {
		return err
	}
	s.server = srv
	return nil
}

// Stop stops the http server
func (s ServerHTTP) Stop(c *cli.Context) error {
	if s.server != nil {
		ctx := context.Background()
		return s.server.Shutdown(ctx)
	}
	// silently return nil if serves has not been launched
	return nil
}
