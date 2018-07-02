package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/skycoin/services/eth-balance-extractor/extractor"
)

func main() {
	o := extractor.NewOrchestrator(1000000, 10)

	go o.StartScanning()

	reader := bufio.NewReader(os.Stdin)
	for {
		reader.ReadString('\n')
	}

	fmt.Println("Shutting down...")
	// e.StopExtraction()
	// scanner.StopScanning()
	// fmt.Println("Last scanned block is: ", e.LastProcessedBlock)
}
