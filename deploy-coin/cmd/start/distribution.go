package main

import (
	"github.com/skycoin/services/deploy-coin/common"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
)

func makeDistributionTx(nc NodeConfig, wc common.GenesisWalletConfig,
	d *daemon.Daemon) (coin.Transaction, error) {

	var tx coin.Transaction

	// Get upnspnets from genesis block
	gb, err := d.Visor.GetSignedBlock(0)
	if err != nil {
		return tx, err
	}
	txIn, err := coin.CreateUnspent(gb.Head, gb.Body.Transactions[0], 0)
	if err != nil {
		return tx, err
	}

	tx.PushInput(txIn.Hash())

	// Create addresses to distribute by inital coin volume
	// First address is address of genesis block, so it is skipped
	// Private key, used to sign genesis block, is used to sign each output
	addrSk := cipher.GenerateDeterministicKeyPairs([]byte(wc.Seed), int(wc.Addresses))
	for _, sk := range addrSk {
		addr := cipher.AddressFromSecKey(sk)
		tx.PushOutput(addr, wc.CoinsPerAddress*1e6, 1)
	}

	tx.SignInputs([]cipher.SecKey{nc.BlockchainSeckey})

	tx.UpdateHeader()

	if err = tx.VerifyInput(coin.UxArray{txIn}); err != nil {
		return tx, err
	}

	if err = tx.Verify(); err != nil {
		return tx, err
	}

	return tx, nil
}
