package btc

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/skycoin/services/otc/pkg/otc"
)

type Client interface {
	GetBlockHash(int64) (*chainhash.Hash, error)
	GetBlockVerboseTx(*chainhash.Hash) (*btcjson.GetBlockVerboseResult, error)
	GetBlockCount() (int64, error)
	WaitForShutdown()
}

type Connection struct {
	Logs    *log.Logger
	Client  Client
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
		Logs:    log.New(os.Stdout, "", log.LstdFlags),
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
						c.Logs.Printf("waiting for block: %d\n", height)
					} else {
						c.Logs.Printf("scan error: %v\n", err)
					}

					time.Sleep(time.Minute)
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
