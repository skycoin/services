package multi

import (
	"crypto/rand"
	"errors"
	"fmt"

	"encoding/hex"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/gui"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

var getBalanceAddresses = func(client ClientApi, addresses []string) (*wallet.BalancePair, error) {
	return client.Balance(addresses)
}

// ClientApi describes skycoin client API
type ClientApi interface {
	Transaction(string) (*visor.TransactionResult, error)
	InjectTransaction(string) (string, error)
	Balance([]string) (*wallet.BalancePair, error)
}

// SkyСoinService provides generic access to various coins API
type SkyСoinService struct {
	client       ClientApi
	checkBalance func(client ClientApi, addresses []string) (*wallet.BalancePair, error)
}

// NewSkyService returns new multicoin generic service
func NewSkyService(n *Node) *SkyСoinService {
	s := &SkyСoinService{
		client:       gui.NewClient(fmt.Sprintf("%s:%d", n.GetNodeHost(), n.GetNodePort())),
		checkBalance: getBalanceAddresses,
	}
	return s
}

func getRand() []byte {
	return cipher.RandByte(1024)
}

// GenerateKeyPair generates key pairs
func (s *SkyСoinService) GenerateKeyPair() *KeysResponse {
	seed := getRand()
	rand.Read(seed)
	pub, sec := cipher.GenerateDeterministicKeyPair(seed)
	return &KeysResponse{
		Private: sec.Hex(),
		Public:  pub.Hex(),
	}
}

// GenerateAddr generates rawAddress, private keys, pubkeys from deterministic seed
func (s *SkyСoinService) GenerateAddr(pubStr string) (maddr *AddressResponse, err error) {
	maddr = &AddressResponse{}

	pubKey, err := cipher.PubKeyFromHex(pubStr)

	if err != nil {
		return nil, err
	}

	address := cipher.AddressFromPubKey(pubKey)

	maddr.Address = address.String()
	return maddr, nil
}

// CheckBalance check the balance (and get unspent outputs) for an rawAddress
func (s *SkyСoinService) CheckBalance(addr string) (*BalanceResponse, error) {
	addressesToGetBalance := make([]string, 0, 1)
	addressesToGetBalance = append(addressesToGetBalance, addr)
	balanceResult, err := s.checkBalance(s.client, addressesToGetBalance)

	if err != nil {
		return nil, err
	}

	return &BalanceResponse{
		Address: addr,
		Hours:   balanceResult.Confirmed.Hours,
		Balance: balanceResult.Confirmed.Coins,
	}, nil
}

// SignTransaction sign a raw transaction with provided private key
func (s *SkyСoinService) SignTransaction(secKey, rawTx string) (response *TransactionSignResponse, err error) {
	response = &TransactionSignResponse{}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error signing transaction %s", r)
		}
	}()

	b, err := hex.DecodeString(rawTx)
	if err != nil {
		fmt.Printf("invalid raw transaction: %v\n", err)
		return nil, err
	}

	tx, err := coin.TransactionDeserialize(b)

	cipherSecKey, err := cipher.SecKeyFromHex(secKey)
	if err != nil {
		return nil, err
	}

	keys := []cipher.SecKey{cipherSecKey}
	tx.SignInputs(keys)
	tx.UpdateHeader()

	// Return raw transaction(hex of signed transaction)
	response.Transaction = tx.TxIDHex()

	if len(tx.Sigs) > 0 {
		response.Signid = tx.Sigs[0].Hex()
	}

	return response, nil
}

// CheckTransactionStatus check the status of a transaction (tracks transactions by transaction hash)
func (s *SkyСoinService) CheckTransactionStatus(txID string) (*visor.TransactionStatus, error) {
	// validate the txid
	_, err := cipher.SHA256FromHex(txID)
	if err != nil {
		return nil, errors.New("invalid txid")
	}
	tx, err := s.client.Transaction(txID)

	if err != nil {
		return nil, err
	}

	return &tx.Status, nil
}

// InjectTransaction send transaction to the network to be included
// in list of unconfirmed transactions.
func (s *SkyСoinService) InjectTransaction(rawtx string) (*Transaction, error) {
	injectedTx, err := s.client.InjectTransaction(rawtx)
	if err != nil {
		return nil, err
	}
	tx, err := s.client.Transaction(injectedTx)
	if err != nil {
		return nil, err
	}

	var tStatus string

	if tx.Status.Confirmed {
		tStatus = "confirmed"
	} else {
		tStatus = "unconfirmed"
	}

	return &Transaction{
		Transid: injectedTx,
		Status:  tStatus,
	}, nil
}
