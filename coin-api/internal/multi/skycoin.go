package multi

import (
	"crypto/rand"
	"fmt"

	"github.com/skycoin/services/coin-api/internal/model"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
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

// GenerateAddr generates address, private keys, pubkeys from deterministic seed
func (s *GenericСoinService) GenerateAddr(count int, hideSecret bool) (*model.Response, error) {
	seed := cipher.SumSHA256(cipher.RandByte(1024)).Hex()
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
		Status: model.ResultOk,
		Code:   0,
		Result: &model.AddressResponse{
			Address: adrss[0].String(),
		},
	}

	return &rsp, nil
}

// GenerateKeyPair generates key pairs
func (s *GenericСoinService) GenerateKeyPair() *model.Response {
	seed := make([]byte, 256)
	rand.Read(seed)
	pub, sec := cipher.GenerateDeterministicKeyPair(seed)
	// address := cipher.AddressFromSecKey(sec)
	// responseParams := map[string]interface{}{
	// 	"publicKey": pub.Hex(),
	// 	"secretKey": sec.Hex(),
	// 	"address":   address.String(),
	// }
	rsp := model.Response{
		Status: model.ResultOk,
		Code:   0,
		Result: &model.KeysResponse{
			Private: pub.Hex(),
			Public:  sec.Hex(),
		},
	}

	// spew.Dump(rsp)
	return &rsp
}

// CheckBalance check the balance (and get unspent outputs) for an address
func (s *GenericСoinService) CheckBalance() {

}

// SignTransaction sign a transaction
func (s *GenericСoinService) SignTransaction() {

}

// CheckTransactionStatus check the status of a transaction (tracks transactions by transaction hash)
func (s *GenericСoinService) CheckTransactionStatus() {

}

// InjectTransaction inject transaction into network
func (s *GenericСoinService) InjectTransaction() {

}
