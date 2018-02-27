package multi

import (
	"crypto/rand"
	"fmt"

	"bytes"
	"strconv"

	"github.com/skycoin/services/coin-api/internal/model"
	"github.com/skycoin/skycoin/src/api/cli"
	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
	// gcli "github.com/urfave/cli"
)

// GenericСoinService provides generic access to various coins API
type GenericСoinService struct {
	// client interface{} // coin client API
}

// NewMultiCoinService returns new multicoin generic service
func NewMultiCoinService() *GenericСoinService {
	//TODO: implement skycoin here
	// connect to skycoin somehow
	// wallet.CreateAddresses()
	return &GenericСoinService{}
}

func getRand() []byte {
	return cipher.RandByte(1024)
}
func getSeed() string {
	return cipher.SumSHA256(getRand()).Hex()
}

// GenerateAddr generates address, private keys, pubkeys from deterministic seed
func (s *GenericСoinService) GenerateAddr(count int, hideSecret bool) (*model.Response, error) {
	seed := getSeed()
	w, err := wallet.CreateAddresses(wallet.CoinTypeSkycoin, seed, count, hideSecret)
	if err != nil {
		return nil, err
	}
	wl, err := w.ToWallet()
	adrss := wl.GetAddresses()
	if len(adrss) == 0 {
		return nil, fmt.Errorf("Unable to get wallet address, number of wallets is %d", len(adrss))
	}
	rsp := model.Response{
		Status: model.StatusOk,
		Code:   0,
		Result: &model.AddressResponse{
			Address: adrss[0].String(),
		},
	}

	return &rsp, nil
}

// GenerateKeyPair generates key pairs
func (s *GenericСoinService) GenerateKeyPair() *model.Response {
	seed := getRand()
	rand.Read(seed)
	pub, sec := cipher.GenerateDeterministicKeyPair(seed)
	// address := cipher.AddressFromSecKey(sec)
	rsp := model.Response{
		Status: model.StatusOk,
		Code:   0,
		Result: &model.KeysResponse{
			Private: pub.Hex(),
			Public:  sec.Hex(),
		},
	}

	return &rsp
}

// CheckBalance check the balance (and get unspent outputs) for an address
func (s *GenericСoinService) CheckBalance(wltFile string, addr int) (*model.Response, error) {
	// wallet.LoadWallets(wltsDir)
	//TODO: probably i have to just get unspent outputs?
	wlt, err := wallet.Load(wltFile)
	if err != nil {
		return nil, err
	}
	addresses := wlt.GetAddresses()

	addressesToGetBalance := make([]string, 0, 1)
	for addWlt := range addresses {
		if addWlt == addr {
			addressesToGetBalance = append(addressesToGetBalance, strconv.Itoa(addr))
		}
	}

	webrpcClient := &webrpc.Client{
		Addr: "someaddr",
		//TODO: fill this rpc client by data including "someaddr"
	}
	balanceResult, err := cli.GetBalanceOfAddresses(webrpcClient, addressesToGetBalance)
	if err != nil {
		return nil, err
	}

	rsp := model.Response{
		Status: model.StatusOk,
		Code:   model.CodeNoError,
		Result: &model.BalanceResponse{
			Address: strconv.Itoa(addr),
			Balance: balanceResult.Spendable.Coins,
			Coin:    model.Coin{
			//TODO: fill data here
			},
		},
	}
	return &rsp, nil
}

// SignTransaction sign a transaction
func (s *GenericСoinService) SignTransaction(transid string) (*model.Response, error) {
	//TODO: VERIFY this sign transaction logic
	var buf bytes.Buffer
	buf.WriteString(transid)
	strbytes := buf.Bytes()
	var secKey cipher.SecKey
	rsp := &model.Response{}
	defer func() {
		if r := recover(); r != nil {
			rsp.Status = model.StatusError
			rsp.Code = -124
			rsp.Result = &model.TransactionSign{}
		}
	}()
	secKey = cipher.NewSecKey(strbytes)
	// coin.NewBlock()
	trans := coin.Transaction{
	//TODO: some creds here?
	}
	keysSec := make([]cipher.SecKey, 0, 1)
	keysSec = append(keysSec, secKey)
	trans.SignInputs(keysSec)
	//TODO: maybe we have to show all signatures?
	signid := trans.Sigs[0]
	rsp.Status = model.StatusOk
	rsp.Code = 0
	rsp.Result = &model.TransactionSign{
		Signid: signid.Hex(),
	}
	return rsp, nil
}

// CheckTransactionStatus check the status of a transaction (tracks transactions by transaction hash)
func (s *GenericСoinService) CheckTransactionStatus() {

}

// InjectTransaction inject transaction into network
func (s *GenericСoinService) InjectTransaction(rawtx string) (*model.Response, error) {
	cli := &webrpc.Client{
	//TODO: insert credentials here
	}
	injectedT, err := cli.InjectTransactionString(rawtx)
	if err != nil {
		return nil, err
	}
	statusT, err := cli.GetTransactionByID(injectedT)
	if err != nil {
		return nil, err
	}

	var tStatus string
	if statusT.Transaction.Status.Confirmed {
		tStatus = "confirmed"
	} else {
		tStatus = "unconfirmed"
	}
	rsp := model.Response{
		Status: model.StatusOk,
		Code:   0,
		Result: &model.Transaction{
			Transid: injectedT,
			Status:  tStatus,
		},
	}
	return &rsp, nil
}
