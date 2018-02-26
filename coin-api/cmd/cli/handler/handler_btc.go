package handler

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"log"
)

const (
	minBtcAddrLen = 26
	maxBtcAddrLen = 35
)

// BTC is a cli bitcoin handler
type BTC struct{}

// NewBTC returns new bitcoin handler instance
func NewBTC() *BTC {
	return &BTC{}
}

// GenerateAddress generates addresses and keypairs for bitcoin
func (b *BTC) GenerateAddress(c *cli.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	publicKey := c.Args().Get(1)

	params := map[string]interface{}{
		"publicKey": publicKey,
	}

	resp, err := doRequest("generateAddr", params)
	if err != nil {
		return err
	}
	log.Printf("Address %s created\n", resp)

	return nil
}

// GenerateKeyPair generates keypair for bitcoin
func (b *BTC) GenerateKeyPair(c *cli.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	resp, err := doRequest("generateKeyPair", nil)
	if err != nil {
		return err
	}
	log.Printf("Key %s created\n", resp)
	return nil
}

// CheckBalance checks bitcoin balance
func (b *BTC) CheckBalance(c *cli.Context) error {
	//TODO: get request info, call appropriate handler from internal btc, don't pass echo context further
	// deal with io.Reader interface
	addr := c.Args().First()

	if len(addr) > 35 || len(addr) < 26 {
		err := errors.New(fmt.Sprintf("Address lenght must be between %d and %d",
			minBtcAddrLen, maxBtcAddrLen))
		return err
	}

	params := map[string]interface{}{
		"address": addr,
	}

	resp, err := doRequest("checkBalance", params)
	if err != nil {
		return err
	}
	log.Printf("Check balance success %s\n", resp)
	return nil
}
