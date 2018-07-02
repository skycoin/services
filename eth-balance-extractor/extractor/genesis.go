package extractor

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core"
)

func ProcessGenesisBlock() {

	g := core.DefaultGenesisBlock()

	// ethInstance := eth.Ethereum{}
	// blockchain := ethInstance.BlockChain()
	// genesis := blockchain.Genesis()
	b := g.ToBlock(nil)
	fmt.Println(b)
}
