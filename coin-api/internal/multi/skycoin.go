package multi

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/gui"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

var getBalanceAddresses = func(client ClientApi, addresses []string) (*wallet.BalancePair, error) {
	return client.Balance(addresses)
}

// ClientApi describes skycoin RPC client API
type ClientApi interface {
	Transaction(string) (*visor.TransactionResult, error)
	InjectTransaction(string) (string, error)
	Balance([]string) (*wallet.BalancePair, error)
}

// SkyСoinService provides generic access to various coins API
type SkyСoinService struct {
	client ClientApi
	// client       *webrpc.Client
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
		Private: pub.Hex(),
		Public:  sec.Hex(),
	}
}

// GenerateAddr generates rawAddress, private keys, pubkeys from deterministic seed
func (s *SkyСoinService) GenerateAddr(pubStr string) (maddr *AddressResponse, err error) {
	maddr = &AddressResponse{}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error generating rawAddress %s", r)
		}
	}()
	pubKey := cipher.MustPubKeyFromHex(pubStr)
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

// SignTransaction sign a transaction
func (s *SkyСoinService) SignTransaction(transID, srcTrans string) (rsp *TransactionSign, err error) {
	rsp = &TransactionSign{}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error signing transaction %s", r)
		}
	}()
	cipherSecKey, err := cipher.SecKeyFromHex(transID)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	buf.WriteString(srcTrans)

	ux := &coin.UxBody{
		SrcTransaction: cipher.SumSHA256(buf.Bytes()),
		Address:        cipher.AddressFromSecKey(cipherSecKey),
		// Coins TODO: maybe we have to receive coins here ?
		// Hours TODO: maybe we have to receive hours here ?
	}
	secKeyTrans := []cipher.SecKey{cipherSecKey}
	trans := &coin.Transaction{}
	uxHash := ux.Hash()
	trans.PushInput(uxHash)
	trans.SignInputs(secKeyTrans)
	//TODO: DO I need it here? -> PushOutput Adds a TransactionOutput, sending coins & hours to an Address
	//TODO: maybe we have to show all signatures?
	signid := trans.Sigs[0]
	rsp.Signid = signid.Hex()
	return rsp, nil
}

// CheckTransactionStatus check the status of a transaction (tracks transactions by transaction hash)
func (s *SkyСoinService) CheckTransactionStatus(txID string) (*visor.TransactionStatus, error) {
	// validate the txid
	_, err := cipher.SHA256FromHex(txID)
	if err != nil {
		return nil, errors.New("invalid txid")
	}
	status, err := s.client.Transaction(txID)

	if err != nil {
		return nil, err
	}

	return &status.Status, nil
}

// InjectTransaction inject transaction into network
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
