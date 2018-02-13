package coin_api

import (
	"github.com/shopspring/decimal"
	"github.com/skycoin/skycoin/src/cipher"
)

func GenerateBtcAddr(pubKey cipher.PubKey) string {
	return cipher.BitcoinAddressFromPubkey(pubKey)
}

func GenerateKeyPair() (cipher.PubKey, cipher.SecKey) {
	return cipher.GenerateDeterministicKeyPair([]byte(""))
}

func CheckBalance() (decimal.Decimal, error) {
	return decimal.NewFromString("0.0")
}
