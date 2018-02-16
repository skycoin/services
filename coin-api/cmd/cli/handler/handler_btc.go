package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"time"
)

const (
	clientTimeout = time.Second * 10
	minBtcAddrLen = 26
	maxBtcAddrLen = 35
)

// BTC is a cli bitcoin handler
type BTC struct{}

// NewBTC returns new bitcoin handler instance
func NewBTC() *BTC {
	return &BTC{}
}

// GenerateKeyPair generates keypair for bitcoin
func (b *BTC) GenerateKeyPair(c *cli.Context) error {
	req, err := http.NewRequest(http.MethodPost, "/keys", nil)

	if err != nil {
		return err
	}

	resp, err := http.Client{
		Timeout: clientTimeout,
	}.Do(req)

	if err != nil {
		return err
	}

	log.Printf("Key %s created\n", resp)
	return nil
}

// GenerateAddress generates addresses and keypairs for bitcoin
func (b *BTC) GenerateAddress(c *cli.Context) error {
	publicKey := c.Args().Get(1)

	params := map[string]interface{}{
		"publicKey": publicKey,
	}

	data, err := json.Marshal(params)

	if err != nil {
		return err
	}

	body := bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, "/address", body)

	if err != nil {
		return err
	}

	resp, err := http.Client{
		Timeout: clientTimeout,
	}.Do(req)

	if err != nil {
		return err
	}

	log.Printf("Address %s created\n", resp)

	return nil
}

// CheckBalance checks bitcoin balance
func (b *BTC) CheckBalance(c *cli.Context) error {
	addr := c.Args().First()

	if len(addr) > 35 || len(addr) < 26 {
		err := errors.New(fmt.Sprintf("Address lenght must be between %d and %d",
			minBtcAddrLen, maxBtcAddrLen))
		return err
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/address/%s", addr), nil)

	if err != nil {
		return err
	}

	resp, err := http.Client{
		Timeout: clientTimeout,
	}.Do(req)

	if err != nil {
		return err
	}

	log.Printf("Check balance success %s\n", resp)
	return nil
}
