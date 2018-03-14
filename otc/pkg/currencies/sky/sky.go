package sky

import (
	"fmt"

	"github.com/skycoin/services/otc/pkg/otc"
	"github.com/skycoin/skycoin/src/api/cli"
	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/wallet"
)

type Connection struct {
	Client    *webrpc.Client
	Wallet    *wallet.Wallet
	FromAddrs []string
}

func New(conf *otc.Config) (*Connection, error) {
	c := &webrpc.Client{Addr: conf.SKY.Node}
	if s, err := c.GetStatus(); err != nil {
		return nil, err
	} else if !s.Running {
		return nil, fmt.Errorf("node isn't running at %s", conf.SKY.Node)
	}

	w, err := wallet.NewWallet(
		conf.SKY.Name,
		wallet.Options{
			Coin:  wallet.CoinTypeSkycoin,
			Label: conf.SKY.Name,
			Seed:  conf.SKY.Seed,
		},
	)
	if err != nil {
		return nil, err
	}

	// TODO: config to put coins in one address?
	_ = w.GenerateAddresses(10)
	conn := &Connection{Wallet: w, Client: c}
	conn.FromAddrs = conn.getFromAddrs()

	return conn, nil
}

func (c *Connection) Used() ([]string, error) {
	var (
		addrs = c.Wallet.GetAddresses()
		out   = make([]string, len(addrs), len(addrs))
	)

	// convert to string
	for i := range addrs {
		out[i] = addrs[i].String()
	}

	return out, nil
}

func (c *Connection) Balance(addr string) (uint64, error) {
	unspent, err := c.Client.GetUnspentOutputs([]string{addr})
	if err != nil {
		return 0, err
	}

	balance, err := unspent.Outputs.SpendableOutputs().Balance()
	if err != nil {
		return 0, err
	}

	return balance.Coins, nil
}

func (c *Connection) Holding() (uint64, error) {
	out, err := cli.GetWalletOutputs(c.Client, c.Wallet)
	if err != nil {
		return 0, err
	}

	bal, err := out.Outputs.SpendableOutputs().Balance()
	if err != nil {
		return 0, err
	}

	return bal.Coins, nil
}

func (c *Connection) Confirmed(txid string) (bool, error) {
	tx, err := c.Client.GetTransactionByID(txid)
	if err != nil {
		return false, err
	}

	return tx.Transaction.Status.Confirmed, nil
}

func (c *Connection) getFromAddrs() []string {
	addrs := c.Wallet.GetAddresses()

	if len(addrs) == 0 {
		addrs = c.Wallet.GenerateAddresses(1)
	}

	out := make([]string, len(addrs))
	for i := range addrs {
		out[i] = addrs[i].String()
	}

	return out
}

func (c *Connection) Send(addr string, amount uint64) (string, error) {
	// create sky transaction
	tx, err := cli.CreateRawTx(c.Client, c.Wallet, c.FromAddrs, c.FromAddrs[0],
		[]cli.SendAmount{{Addr: addr, Coins: amount}},
	)
	if err != nil {
		return "", err
	}

	// inject and get txId
	txid, err := c.Client.InjectTransaction(tx)
	if err != nil {
		return "", err
	}

	return txid, nil
}

func (c *Connection) Address() (string, error) {
	addr := c.Wallet.GenerateAddresses(1)
	if addr == nil {
		return "", fmt.Errorf("error generating address")
	}

	return addr[0].String(), nil
}

func (c *Connection) Connected() (bool, error) {
	res, err := c.Client.GetStatus()
	return res.Running, err
}

func (c *Connection) Stop() error { return nil }
