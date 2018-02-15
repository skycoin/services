package multi

// GenericСoinService provides generic access to various coins API
type GenericСoinService struct {
	// client interface{} // coin client API
}

// NewMultiCoinService returns new multicoin generic service
func NewMultiCoinService() *GenericСoinService {
	return &GenericСoinService{}
}

// GenerateAddr generates address, private keys, pubkeys from deterministic seed
func (s *GenericСoinService) GenerateAddr() {

}

// GenerateKeyPair generate private keys, pubkeys from deterministic seed
func (s *GenericСoinService) generateKeyPair() {

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
