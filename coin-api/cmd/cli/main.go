package main

import (
	"os"

	"github.com/skycoin/services/coin-api/cli"
)

func main() {
	err := cli.App.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
