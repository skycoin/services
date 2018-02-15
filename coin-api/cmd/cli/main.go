package main

import (
	"os"

	"github.com/skycoin/services/coin-api/cmd/cli/handler"
	"github.com/skycoin/services/coin-api/cmd/servd"
	"github.com/urfave/cli"
)

var (
	rpcaddr  = new(string)
	endpoint string
)

func main() {
	// App is a cli app
	hBTC := handler.NewBTC()
	hMulti := handler.NewMulti()
	cliapp := cli.App{
		Commands: []cli.Command{
			{
				Name: "btc",
				Subcommands: cli.Commands{
					cli.Command{
						Name:   "generatekeys",
						Usage:  "Generate key pair",
						Action: hBTC.GenerateKeyPair,
					},
					cli.Command{
						Name:      "generateaddr",
						Usage:     "Generate BTC addr",
						ArgsUsage: "<publicKey>",
						Action:    hBTC.GenerateAddress,
					},
					cli.Command{
						Name:      "checkbalance",
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
			{
				Name: "coin",
				Subcommands: cli.Commands{
					cli.Command{
						Name:   "generatekeys",
						Usage:  "Generate key pair",
						Action: hMulti.GenerateKeyPair,
					},
					cli.Command{
						Name:      "generateaddr",
						Usage:     "Generate BTC addr",
						ArgsUsage: "<publicKey>",
						Action:    hMulti.GenerateAddress,
					},
					cli.Command{
						Name:      "checkbalance",
						Usage:     "Check BTC balance",
						ArgsUsage: "<address>",
						Action:    hMulti.CheckBalance,
					},
				},
			},
			{
				Name: "server",
				Subcommands: cli.Commands{
					cli.Command{
						Name:   "start",
						Usage:  "Start HTTP Server",
						Action: servd.Start,
					},
					cli.Command{
						Name:   "stop",
						Usage:  "Stop HTTP Server",
						Action: servd.Start,
					},
				},
			},
		},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "rpc", Destination: rpcaddr, Value: "localhost:12345"},
		},
	}
	err := cliapp.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
