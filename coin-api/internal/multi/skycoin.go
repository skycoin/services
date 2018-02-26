package multi

import (
	"crypto/rand"

	"github.com/skycoin/services/coin-api/internal/model"
	"github.com/skycoin/skycoin/src/cipher"
)

// GenericСoinService provides generic access to various coins API
type GenericСoinService struct {
	// client interface{} // coin client API
}

// NewMultiCoinService returns new multicoin generic service
func NewMultiCoinService() *GenericСoinService {
	return &GenericСoinService{}
}

// GenerateAddr generates address, private keys, pubkeys from deterministic seed
func (s *GenericСoinService) GenerateAddr() *model.Response {
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
