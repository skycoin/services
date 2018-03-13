package btc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/skycoin/services/otc/pkg/otc"
)

const (
	EXPLORER_PATH         = "https://blockchain.info/unspent?active="
	EXPLORER_PATH_TESTNET = "https://testnet.blockchain.info/unspent?active="
)

type Connection struct {
	Client  *rpcclient.Client
	Account string
	Testnet bool
}

func New(conf *otc.Config) (*Connection, error) {
	// get tls certs for websocket connection
	certs, err := ioutil.ReadFile(
		filepath.Join(
			btcutil.AppDataDir("btcwallet", false),
			"rpc.cert",
		),
	)
	if err != nil {
		return nil, err
	}

	// connect to btc node
	client, err := rpcclient.New(
		&rpcclient.ConnConfig{
			Host:         conf.BTC.Node,
			HTTPPostMode: true,
			User:         conf.BTC.User,
			Pass:         conf.BTC.Pass,
			Certificates: certs,
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Connection{
		Client:  client,
		Account: conf.BTC.Account,
		Testnet: conf.BTC.Testnet,
	}, nil
}

// TODO
func (c *Connection) Holding() (uint64, error) {
	return 0, nil
}

func (c *Connection) Balance(addr string) (uint64, error) {
	var path string
	if c.Testnet {
		path = EXPLORER_PATH_TESTNET
	} else {
		path = EXPLORER_PATH
	}

	resp, err := http.Get(path + addr + "&confirmations=1")
	if err != nil {
		return 0, err
	}

	var data struct {
		UnspentOutputs []struct {
			Value uint64 `json:"value"`
		} `json:"unspent_outputs"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, nil
	}

	var sum uint64
	for _, o := range data.UnspentOutputs {
		sum += o.Value
	}

	return sum, nil
}

func (c *Connection) Confirmed(txid string) (bool, error) {
	return false, nil
}

func (c *Connection) Send(addr string, amount uint64) (string, error) {
	return "", nil
}

func (c *Connection) Address() (string, error) {
	addr, err := c.Client.GetNewAddress(c.Account)
	if err != nil {
		return "", err
	}

	return addr.EncodeAddress(), nil
}

func (c *Connection) Connected() (bool, error) {
	return !c.Client.Disconnected(), c.Client.Ping()
}

func (c *Connection) Stop() error {
	c.Client.Shutdown()
	c.Client.WaitForShutdown()
	return nil
}
