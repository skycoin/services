package main

import (
	"github.com/mihis/services/deploy-coin/common"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
)

func makeDistributionTx(nc NodeConfig, dc common.DistributionConfig,
	d *daemon.Daemon) (coin.Transaction, error) {

	gb, err := d.Visor.GetBlock(0)
	if err != nil {
		return coin.Transaction{}, err
	}

	txIn, err := coin.CreateUnspent(gb.Head, gb.Body.Transactions[0], 0)
	if err != nil {
		return coin.Transaction{}, err
	}

	var tx coin.Transaction
	tx.PushInput(txIn.Hash())

	for i := range dc.Addresses {
		addr := cipher.MustDecodeBase58Address(dc.Addresses[i])
		tx.PushOutput(addr, dc.CoinsPerAddress, 1)
	}

	keys := []cipher.SecKey{
		nc.BlockchainSeckey,
	}
	tx.SignInputs(keys)

	tx.UpdateHeader()

	if err = tx.VerifyInput(coin.UxArray{txIn}); err != nil {
		return tx, err
	}

	if err = tx.Verify(); err != nil {
		return tx, err
	}

	return tx, nil
}
