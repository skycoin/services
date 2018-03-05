// NOTE: this package is not to be used outside of the load-testing utility.
// Most of these functions are rough copies from the skycoin/src/api/cli
// package with changes made to create transactions with zero CoinHours on
// the outputs.
package staggering

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/api/cli"
	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

var (
	// ErrTemporaryInsufficientBalance is returned if a wallet does not have enough balance for a spend, but will have enough after unconfirmed transactions confirm
	ErrTemporaryInsufficientBalance = errors.New("balance is not sufficient. Balance will be sufficient after unconfirmed transactions confirm")

	ErrAddress = errors.New("bad address")
)

// CreateRawTx creates a transaction from a set of addresses contained in a loaded *wallet.Wallet
func CreateRawTx(c *webrpc.Client, wlt *wallet.Wallet, inAddrs []string, chgAddr string, toAddrs []cli.SendAmount) (*coin.Transaction, error) {
	if err := validateSendAmounts(toAddrs); err != nil {
		return nil, err
	}

	// Get unspent outputs of those addresses
	unspents, err := c.GetUnspentOutputs(inAddrs)
	if err != nil {
		return nil, err
	}

	inUxs, err := unspents.Outputs.SpendableOutputs().ToUxArray()
	if err != nil {
		return nil, err
	}

	txn, err := createRawTx(unspents.Outputs, wlt, inAddrs, chgAddr, toAddrs)
	if err != nil {
		return nil, err
	}

	// filter out unspents which are not used in transaction
	var inUxsFiltered coin.UxArray
	for _, h := range txn.In {
		for _, u := range inUxs {
			if h == u.Hash() {
				inUxsFiltered = append(inUxsFiltered, u)
			}
		}
	}

	// TODO -- remove me -- reimplementation of visor.VerifySingleTxnSoftConstraints minus
	// the parts that require block head data, which is not available from the RPC API (see below)
	if err := verifyTransactionConstraints(txn, inUxsFiltered, visor.DefaultMaxBlockSize); err != nil {
		return nil, err
	}

	// TODO -- verify against soft and hard constraints
	// Need to get the head block to do verification.
	// The head block is not exposed over the JSON RPC, which webrpc.Client uses.
	// Need to remove the JSON RPC API and have the client make requests to the HTTP API.
	// Once the HTTP API is used,
	// Need to request /blockchain/metadata to get the head block time
	// This could lead to race conditions; /blockchain/metadata should return the full head, or have an API endpoint
	// just for the head, and/or include the head block in the get_outputs response
	// The head block is used for calculating inUxs's coin hours.
	// if err := visor.VerifySingleTxnSoftConstraints(txn, inUxs, visor.DefaultMaxBlockSize); err != nil {
	//     return nil, err
	// }
	// if err := visor.VerifySingleTxnHardConstraints(txn, head, inUxs); err != nil {
	// 	return nil, err
	// }

	return txn, nil
}

func validateSendAmounts(toAddrs []cli.SendAmount) error {
	for _, arg := range toAddrs {
		// validate to address
		_, err := cipher.DecodeBase58Address(arg.Addr)
		if err != nil {
			return ErrAddress
		}

		if arg.Coins == 0 {
			return errors.New("Cannot send 0 coins")
		}
	}

	if len(toAddrs) == 0 {
		return errors.New("No destination addresses")
	}

	return nil
}

// TODO -- remove me -- reimplementation of visor.VerifySingleTxnSoftConstraints and HardConstraints
// minus the parts that require block head data, which is not available from the RPC API (see below)
func verifyTransactionConstraints(txn *coin.Transaction, uxIn coin.UxArray, maxSize int) error {
	// SOFT constraints:

	if txn.Size() > maxSize {
		return errors.New("Transaction size bigger than max block size")
	}

	if visor.TransactionIsLocked(uxIn) {
		return errors.New("Transaction has locked address inputs")
	}

	// Ignore transactions that do not conform to decimal restrictions
	for _, o := range txn.Out {
		if err := visor.DropletPrecisionCheck(o.Coins); err != nil {
			return err
		}
	}

	// HARD constraints:

	if err := txn.Verify(); err != nil {
		return err
	}

	// Checks whether ux inputs exist,
	// Check that signatures are allowed to spend inputs
	if err := txn.VerifyInput(uxIn); err != nil {
		return err
	}

	// Verify CoinHours do not overflow
	if _, err := txn.OutputHours(); err != nil {
		return err
	}

	// Check that no coins are created or destroyed
	// TODO -- use the correct block head, once we have it from the API
	// For now it doesn't matter, the block head is used to calculate the uxOut hours,
	// but we're not validating the hours
	uxOut := coin.CreateUnspents(coin.BlockHeader{
		BkSeq: 1,
	}, *txn)
	return coin.VerifyTransactionCoinsSpending(uxIn, uxOut)

	// TODO -- use coin.VerifyTransactionHoursSpending, once we have the head block
	// return coin.VerifyTransactionHoursSpending(head.Time(), uxIn, uxOut)
}

func createRawTx(uxouts visor.ReadableOutputSet, wlt *wallet.Wallet, inAddrs []string, chgAddr string, toAddrs []cli.SendAmount) (*coin.Transaction, error) {
	// Calculate total required coins
	var totalCoins uint64
	for _, arg := range toAddrs {
		totalCoins += arg.Coins
	}

	spendOutputs, err := chooseSpends(uxouts, totalCoins)
	if err != nil {
		return nil, err
	}

	keys, err := getKeys(wlt, spendOutputs)
	if err != nil {
		return nil, err
	}

	txOuts, err := makeChangeOut(spendOutputs, chgAddr, toAddrs)
	if err != nil {
		return nil, err
	}

	tx := cli.NewTransaction(spendOutputs, keys, txOuts)

	return tx, nil
}

func makeChangeOut(outs []wallet.UxBalance, chgAddr string, toAddrs []cli.SendAmount) ([]coin.TransactionOutput, error) {
	var totalInCoins, totalInHours, totalOutCoins uint64

	for _, o := range outs {
		totalInCoins += o.Coins
		totalInHours += o.Hours
	}

	if totalInHours == 0 {
		return nil, fee.ErrTxnNoFee
	}

	for _, to := range toAddrs {
		totalOutCoins += to.Coins
	}

	if totalInCoins < totalOutCoins {
		return nil, wallet.ErrInsufficientBalance
	}

	outAddrs := []coin.TransactionOutput{}
	changeAmount := totalInCoins - totalOutCoins

	haveChange := changeAmount > 0
	var totalOutHours uint64 = 0

	if err := fee.VerifyTransactionFeeForHours(totalOutHours, totalInHours-totalOutHours); err != nil {
		return nil, err
	}

	if haveChange {
		outAddrs = append(outAddrs, mustMakeUtxoOutput(chgAddr, changeAmount, totalInHours/2))
	}

	for _, to := range toAddrs {
		outAddrs = append(outAddrs, mustMakeUtxoOutput(to.Addr, to.Coins, 0))
	}

	return outAddrs, nil
}

func chooseSpends(uxouts visor.ReadableOutputSet, coins uint64) ([]wallet.UxBalance, error) {
	// Convert spendable unspent outputs to []wallet.UxBalance
	spendableOutputs, err := visor.ReadableOutputsToUxBalances(uxouts.SpendableOutputs())
	if err != nil {
		return nil, err
	}

	// Choose which unspent outputs to spend
	// Use the MinimizeUxOuts strategy, since this is most likely used by
	// application that may need to send frequently.
	// Using fewer UxOuts will leave more available for other transactions,
	// instead of waiting for confirmation.
	outs, err := wallet.ChooseSpendsMinimizeUxOuts(spendableOutputs, coins)
	if err != nil {
		// If there is not enough balance in the spendable outputs,
		// see if there is enough balance when including incoming outputs
		if err == wallet.ErrInsufficientBalance {
			expectedOutputs, otherErr := visor.ReadableOutputsToUxBalances(uxouts.ExpectedOutputs())
			if otherErr != nil {
				return nil, otherErr
			}

			if _, otherErr := wallet.ChooseSpendsMinimizeUxOuts(expectedOutputs, coins); otherErr != nil {
				return nil, err
			}

			return nil, ErrTemporaryInsufficientBalance
		}

		return nil, err
	}

	return outs, nil
}

func getKeys(wlt *wallet.Wallet, outs []wallet.UxBalance) ([]cipher.SecKey, error) {
	keys := make([]cipher.SecKey, len(outs))
	for i, o := range outs {
		entry, ok := wlt.GetEntry(o.Address)
		if !ok {
			return nil, fmt.Errorf("%v is not in wallet", o.Address.String())
		}

		keys[i] = entry.Secret
	}
	return keys, nil
}

func mustMakeUtxoOutput(addr string, coins, hours uint64) coin.TransactionOutput {
	uo := coin.TransactionOutput{}
	uo.Address = cipher.MustDecodeBase58Address(addr)
	uo.Coins = coins
	uo.Hours = hours
	return uo
}
