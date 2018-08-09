package main

import (
	"fmt"
	"os"

	"github.com/skycoin/services/eth-balance-extractor/extractor"
)

func main() {
	app := extractor.NewApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
