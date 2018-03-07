package skycoin

import (
	"fmt"

	"github.com/skycoin/services/otc/types"
	"github.com/skycoin/skycoin/src/api/cli"
	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/wallet"
)

type Connection struct {
	Wallet *wallet.Wallet
	Client *webrpc.Client
}

func NewConnection(config *types.Config) (*Connection, error) {
	c := &webrpc.Client{Addr: config.Skycoin.Node}
	if s, err := c.GetStatus(); err != nil {
		return nil, err
	} else if !s.Running {
		return nil, fmt.Errorf("node isn't running at %s", config.Skycoin.Node)
	}

	w, err := wallet.NewWallet(
		config.Skycoin.Name,
		wallet.Options{
			Coin:  wallet.CoinTypeSkycoin,
			Label: config.Skycoin.Name,
			Seed:  config.Skycoin.Seed,
		},
	)
	if err != nil {
		return nil, err
	}

	// TODO: config to put coins in one address?
	_ = w.GenerateAddresses(100)

	return &Connection{Wallet: w, Client: c}, nil
}

func (c *Connection) Balance() (uint64, error) {
	out, err := cli.GetWalletOutputs(c.Client, c.Wallet)
	if err != nil {
		return 0, err
	}

	bal, err := out.Outputs.SpendableOutputs().Balance()
	if err != nil {
		return 0, err
	}

	return bal.Coins, nil

	/*
		what a nightmare

		outs, err := out.Outputs.SpendableOutputs().ToUxArray()
		if err != nil {
			return 0, err
		}

		var bal uint64
		for _, o := range outs {
			bal += o.Body.Coins
		}

		return bal, nil
	*/
}
