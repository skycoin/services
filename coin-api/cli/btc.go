package cli

import (
	"github.com/urfave/cli"
	"log"
)

func btcCMD() cli.Command {
	return cli.Command{
		Name:  "generateAddr",
		Usage: "Generate BTC addr",
		ArgsUsage: "<publicKey>",
		Subcommands: cli.Commands{
			generateAddrCMD(),
			generateKeyPairCMD(),
			checkBalanceCMD(),
		},
	}

}


func generateAddrCMD() cli.Command {
	return cli.Command{
		Name:  "generateAddr",
		Usage: "Generate BTC addr",
		ArgsUsage: "<publicKey>",
		Action: func(c *cli.Context) error {
			publicKey := c.Args().Get(1)

			params := map[string]interface{}{
				"publicKey": publicKey,
			}

			resp, err := rpcRequest("generateAddr", params)
			if err != nil {
				return err
			}
			log.Printf("Address %s created\n", resp)
			return nil
		},
	}
}

func generateKeyPairCMD() cli.Command {
	return cli.Command{
		Name:  "generateKeyPair",
		Usage: "Generate key pair",
		Action: func(c *cli.Context) error {
			resp, err := rpcRequest("generateKeyPair", nil)
			if err != nil {
				return err
			}
			log.Printf("Key %s created\n", resp)
			return nil
		},
	}
}

func checkBalanceCMD() cli.Command {
	return cli.Command{
		Name:  "check",
		Usage: "Check BTC balance",
		ArgsUsage: "<address>",
		Action: func(c *cli.Context) error {
			addr := c.Args().Get(1)

			params := map[string]interface{}{
				"addr": addr,
			}

			resp, err := rpcRequest("checkBalance", params)
			if err != nil {
				return err
			}
			log.Printf("Check balance success %s\n", resp)
			return nil
		},
	}
}
