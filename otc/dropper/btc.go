package dropper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/skycoin/services/otc/exchange"
	"github.com/skycoin/services/otc/types"
)

const (
	EXPLORER              = "explorer"
	EXPLORER_PATH         = "https://blockchain.info/unspent?active="
	EXPLORER_PATH_TESTNET = "https://testnet.blockchain.info/unspent?active="

	NODE = "node"
)

type BTCConnection struct {
	client  *rpcclient.Client
	account string
	testnet bool
	source  string
}

func NewBTCConnection(config *types.Config) (*BTCConnection, error) {
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
			Host:         config.Dropper.BTC.Node,
			Endpoint:     "ws",
			User:         config.Dropper.BTC.User,
			Pass:         config.Dropper.BTC.Pass,
			Certificates: certs,
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	// get list of all accounts
	accounts, err := client.ListAccounts()
	if err != nil {
		return nil, err
	}

	// check if at least one is for teller-json
	exists := false
	for account, _ := range accounts {
		if account == config.Dropper.BTC.Account {
			exists = true
		}
	}

	// account for teller doesn't exist, need to create a new one
	if !exists {
		// authenticate with the wallet passphrase
		err = client.WalletPassphrase(config.Dropper.BTC.Account, 2)
		if err != nil {
			return nil, err
		}

		// create new account for generating addresses
		err = client.CreateNewAccount(config.Dropper.BTC.Account)
		if err != nil {
			return nil, err
		}
	}

	return &BTCConnection{
		client:  client,
		account: config.Dropper.BTC.Account,
		testnet: config.Dropper.BTC.Testnet,
		source:  config.Dropper.BTC.Source,
	}, nil
}

func (c *BTCConnection) Generate() (types.Drop, error) {
	addr, err := c.client.GetNewAddress(c.account)
	if err != nil {
		return "", err
	}

	return types.Drop(addr.EncodeAddress()), nil
}

func (c *BTCConnection) Send(drop types.Drop, amount uint64) (string, error) {
	// convert string to btc address
	addr, err := btcutil.DecodeAddress(string(drop), &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}

	// unlock wallet for sending
	if err = c.client.WalletPassphrase(c.account, 2); err != nil {
		return "", err
	}

	// send and get transaction id
	hash, err := c.client.SendToAddress(addr, btcutil.Amount(amount))
	if err != nil {
		return "", err
	}

	return hash.String(), nil
}

func (c *BTCConnection) Value() (uint64, error) {
	value, err := exchange.GetBTCValue()
	return value, err
}

type ExplorerResponse struct {
	UnspentOutputs []UnspentOutput `json:"unspent_outputs"`
}

type UnspentOutput struct {
	Value uint64 `json:"value"`
}

func (c *BTCConnection) Balance(drop types.Drop) (uint64, error) {
	if c.source == EXPLORER {
		var path string
		if c.testnet {
			path = EXPLORER_PATH_TESTNET
		} else {
			path = EXPLORER_PATH
		}

		resp, err := http.Get(path + string(drop) + "&confirmations=1")
		if err != nil {
			return 0, err
		}

		var eo ExplorerResponse
		if err = json.NewDecoder(resp.Body).Decode(&eo); err != nil {
			return 0, nil
		}

		var sum uint64
		for _, o := range eo.UnspentOutputs {
			sum += o.Value
		}

		return sum, nil
	}

	var params *chaincfg.Params
	if c.testnet {
		params = &chaincfg.TestNet3Params
	} else {
		params = &chaincfg.MainNetParams
	}

	// convert address string to btc struct
	addr, err := btcutil.DecodeAddress(string(drop), params)
	if err != nil {
		return 0, nil
	}

	// get unspents of address
	unspents, err := c.client.ListUnspentMinMaxAddresses(
		1,                       // min confirmations
		999999,                  // max confirmations
		[]btcutil.Address{addr}, // only checking one address
	)
	if err != nil {
		return 0, nil
	}

	var sum uint64
	for _, res := range unspents {
		if res.Spendable {
			sum += uint64(res.Amount * 10e7)
		}
	}

	return sum, nil
}

func (c *BTCConnection) Confirmed(hash string) (bool, error) {
	txHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return false, err
	}

	// get transaction from blockchain
	result, err := c.client.GetTransaction(txHash)
	if err != nil {
		return false, err
	}

	// 6 confirmations to be confirmed
	if result.Confirmations < 6 {
		return false, nil
	}

	// tx confirmed
	return true, nil
}

func (c *BTCConnection) Connected() (bool, error) {
	return !c.client.Disconnected(), c.client.Ping()
}

func (c *BTCConnection) Stop() error {
	c.client.Shutdown()
	c.client.WaitForShutdown()
	return nil
}
