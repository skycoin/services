package main

import (
	"os"

	"github.com/skycoin/services/coin-api/cmd/cli/handler"
	"github.com/urfave/cli"
)

var (
	rpcaddr  = new(string)
	endpoint string
)

func main() {
	// App is a cli app
	hBTC := handler.NewBTC()
	var App = cli.App{
		Commands: []cli.Command{
			{
				Name: "btc",
				Subcommands: cli.Commands{
					cli.Command{
						Name:   "generateKeyPair",
						Usage:  "Generate key pair",
						Action: hBTC.GenerateKeyPair,
					},
					cli.Command{
						Name:      "generateAddr",
						Usage:     "Generate BTC addr",
						ArgsUsage: "<publicKey>",
						Action:    hBTC.GenerateAddress,
					},
					cli.Command{
						Name:      "checkBalance",
						Usage:     "Check BTC balance",
						ArgsUsage: "<address>",
						Action:    hBTC.CheckBalance,
					},
				},
				Before: func(c *cli.Context) error {
					endpoint = "btc"
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "rpc", Destination: rpcaddr, Value: "localhost:12345"},
		},
	}
	err := cli.App.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
