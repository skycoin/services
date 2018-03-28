package btc

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/skycoin/services/otc/pkg/otc"
)

type Connection struct {
	Client  *rpcclient.Client
	Account string
	stop    chan struct{}
}

func New(account, pass, rNode, rUser, rPass string) (*Connection, error) {
	certs, err := ioutil.ReadFile(filepath.Join(
		btcutil.AppDataDir("btcwallet", false), "rpc.cert"))
	if err != nil {
		return nil, err
	}

	// connect to btc node
	client, err := rpcclient.New(
		&rpcclient.ConnConfig{
			// TODO: figure out why websockets hang
			HTTPPostMode: true,
			Host:         rNode,
			User:         rUser,
			Pass:         rPass,
			Certificates: certs,
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	accounts, err := client.ListAccounts()
	if err != nil {
		return nil, err
	}

	found := false
	for acc, _ := range accounts {
		if acc == account {
			found = true
		}
	}

	if !found {
		if err = client.WalletPassphrase(pass, 1); err != nil {
			return nil, err
		}
		if err = client.CreateNewAccount(account); err != nil {
			return nil, err
		}
	}

	return &Connection{
		Client:  client,
		Account: account,
		stop:    make(chan struct{}, 0),
	}, nil
}

func (c *Connection) Scan(from uint64) (chan *otc.Block, error) {
	blocks := make(chan *otc.Block, 0)
	height := from

	go func() {
		for {
			select {
			case <-c.stop:
				return
			default:
				block, err := c.Get(height)
				if err != nil {
					// TODO: use variable from btc* package rather than str
					if err.Error() == "-1: Block number out of range" {
						// TODO: log using logger
						fmt.Printf("waiting for block: %d\n", height)
						// wait a minute before checking for the next block
						time.Sleep(time.Minute)
					} else {
						// TODO: log using logger
					}
				} else {
					// send block to scanner
					blocks <- block
					// next iteration attempt to get next block
					height++
				}
			}
		}
	}()

	return blocks, nil
}

func (c *Connection) Get(height uint64) (*otc.Block, error) {
	// get block hash of block at height
	bh, err := c.Client.GetBlockHash(int64(height))
	if err != nil {
		return nil, err
	}

	// get block with tx data
	bb, err := c.Client.GetBlockVerboseTx(bh)
	if err != nil {
		return nil, err
	}

	block := &otc.Block{
		Height:       height,
		Transactions: make(map[string]*otc.Transaction, len(bb.RawTx)),
	}

	for _, tx := range bb.RawTx {
		block.Transactions[tx.Hash] = &otc.Transaction{
			BlockHash:     tx.BlockHash,
			Hash:          tx.Hash,
			Confirmations: tx.Confirmations,
			Out:           make(map[int]*otc.Output, len(tx.Vout)),
		}

		for _, out := range tx.Vout {
			amount, err := btcutil.NewAmount(out.Value)
			if err != nil {
				return nil, err
			}

			block.Transactions[tx.Hash].Out[int(out.N)] = &otc.Output{
				Amount:    uint64(amount),
				Addresses: out.ScriptPubKey.Addresses,
			}
		}
	}

	return block, nil
}

func (c *Connection) Height() (uint64, error) {
	count, err := c.Client.GetBlockCount()
	return uint64(count), err
}

func (c *Connection) Stop() error {
	c.stop <- struct{}{}
	c.Client.WaitForShutdown()
	return nil
}
