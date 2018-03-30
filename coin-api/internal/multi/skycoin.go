package multi

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"reflect"

	"github.com/skycoin/skycoin/src/api/cli"
	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
)

var getBalanceAddresses = func(client WebRPCClientAPI, addresses []string) (*cli.BalanceResult, error) {
	webRPC, ok := client.(*webrpc.Client)
	if !ok {
		panic(fmt.Sprintf("wrong type %s *webrpc.Client expected", reflect.TypeOf(webRPC).String()))
	}
	return cli.GetBalanceOfAddresses(webRPC, addresses)
}

// WebRPCClientAPI describes skycoin RPC client API
type WebRPCClientAPI interface {
	GetTransactionByID(string) (*webrpc.TxnResult, error)
	InjectTransactionString(string) (string, error)
}

// SkyСoinService provides generic access to various coins API
type SkyСoinService struct {
	client WebRPCClientAPI
	// client       *webrpc.Client
	checkBalance func(client WebRPCClientAPI, addresses []string) (*cli.BalanceResult, error)
}

// NewSkyService returns new multicoin generic service
func NewSkyService(n *Node) *SkyСoinService {
	s := &SkyСoinService{
		client: &webrpc.Client{
			Addr: fmt.Sprintf("%s:%d", n.GetNodeHost(), n.GetNodePort()),
		},
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

// GenerateAddr generates address, private keys, pubkeys from deterministic seed
func (s *SkyСoinService) GenerateAddr(pubStr string) (maddr *AddressResponse, err error) {
	maddr = &AddressResponse{}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error generating address %s", r)
		}
	}()
	pubKey := cipher.MustPubKeyFromHex(pubStr)
	address := cipher.AddressFromPubKey(pubKey)

	maddr.Address = address.String()
	return maddr, nil
}

func getBalanceAddress(br *cli.BalanceResult) string {
	if len(br.Addresses) > 0 {
		return br.Addresses[0].Address
	}

	return ""
}

// CheckBalance check the balance (and get unspent outputs) for an address
func (s *SkyСoinService) CheckBalance(addr string) (*BalanceResponse, error) {
	addressesToGetBalance := make([]string, 0, 1)
	addressesToGetBalance = append(addressesToGetBalance, addr)
	balanceResult, err := s.checkBalance(s.client, addressesToGetBalance)
	if err != nil {
		return nil, err
	}

	return &BalanceResponse{
		Address: getBalanceAddress(balanceResult),
		Hours:   balanceResult.Spendable.Hours,
		Balance: balanceResult.Spendable.Coins,
		// balanceResult.
		// 	Coin: model.Coin{
		// //TODO: maybe coin info will required in the nearest future
		// },
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
	status, err := s.client.GetTransactionByID(txID)

	if err != nil {
		return nil, err
	}

	return &status.Transaction.Status, nil
}

// InjectTransaction inject transaction into network
func (s *SkyСoinService) InjectTransaction(rawtx string) (*Transaction, error) {
	injectedT, err := s.client.InjectTransactionString(rawtx)
	if err != nil {
		return nil, err
	}
	statusT, err := s.client.GetTransactionByID(injectedT)
	if err != nil {
		return nil, err
	}

	var tStatus string
	if statusT.Transaction.Status.Confirmed {
		tStatus = "confirmed"
	} else {
		tStatus = "unconfirmed"
	}

	return &Transaction{
		Transid: injectedT,
		Status:  tStatus,
	}, nil
}
