package cli

import "github.com/urfave/cli"

func generateAddrCMD() cli.Command {
	return cli.Command{
		Name:  "generateAddr",
		Usage: "Generate BTC addr",
		Action: func() {
			//TODO(stgleb): Implement
		},
	}
}

func generateKeyPairCMD() cli.Command {
	return cli.Command{
		Name:  "generateKey",
		Usage: "Generate key pair",
		Action: func() {
			//TODO(stgleb): Implement
		},
	}
}

func checkBalanceCMD() cli.Command {
	return cli.Command{
		Name:  "check",
		Usage: "Check BTC balance",
		Action: func() {
			//TODO(stgleb): Implement
		},
	}
}
