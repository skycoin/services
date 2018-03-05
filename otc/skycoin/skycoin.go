package skycoin

import (
	"fmt"

	"github.com/skycoin/services/otc/types"
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

	return &Connection{Wallet: w, Client: c}, nil
}
