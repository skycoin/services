package handler

import (
	"log"

	"github.com/urfave/cli"
)

// BTC is a cli bitcoin handler
type BTC struct{}

// NewBTC returns new bitcoin handler instance
func NewBTC() *BTC {
	return &BTC{}
}

// GenerateAddress generates and keypairs for bitcoin
func (b *BTC) GenerateAddress(c *cli.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
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
}

func (b *BTC) GenerateKeyPair(c *cli.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	resp, err := rpcRequest("generateKeyPair", nil)
	if err != nil {
		return err
	}
	log.Printf("Key %s created\n", resp)
	return nil
}

func (b *BTC) CheckBalance(c *cli.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
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
}
