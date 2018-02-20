package btc

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"

	"encoding/json"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/shopspring/decimal"
	"github.com/skycoin/skycoin/src/cipher"
	"net/http"
	"time"
)

const (
	defaultAddr = "23.92.24.9"
	defaultUser = "YnWD3EmQAOw11IOrUJwWxAThAyobwLC"
	defaultPass = `f*Z"[1215o{qKW{Buj/wheO8@h.}j*u`
	defaultCert = `-----BEGIN CERTIFICATE---—
MIICbTCCAc+gAwIBAgIRAKnAvGj6JobKblRUcmxOqxowCgYIKoZIzj0EAwQwNjEg
MB4GA1UEChMXYnRjZCBhdXRvZ2VuZXJhdGVkIGNlcnQxEjAQBgNVBAMTCWxvY2Fs
aG9zdDAeFw0xNzExMDYwNTMzNDRaFw0yNzExMDUwNTMzNDRaMDYxIDAeBgNVBAoT
F2J0Y2QgYXV0b2dlbmVyYXRlZCBjZXJ0MRIwEAYDVQQDEwlsb2NhbGhvc3QwgZsw
EAYHKoZIzj0CAQYFK4EEACMDgYYABAEYn5Xj5QfV6vK6jjeLnG63H5U8yrga5wYJ
bqBhuSR+540zqVjviZQXDi9OVTcYffDk+VrP2KmD8Q8FW2yFAjo2ewA63DHQibtJ
Jb2bSCSJnMa7MqWeYle61oIwt9wIiq+9gjVIagnlEAOVm86TBeuiCgUu5t3k1CrI
R4XFVPAgDQXnzqN7MHkwDgYDVR0PAQH/BAQDAgKkMA8GA1UdEwEB/wQFMAMBAf8w
VgYDVR0RBE8wTYIJbG9jYWxob3N0hwR/AAABhxAAAAAAAAAAAAAAAAAAAAABhwQX
XBgJhxAmADwBAAAAAPA8kf/+zLGFhxD+gAAAAAAAAPA8kf/+zLGFMAoGCCqGSM49
BAMEA4GLADCBhwJCATk6kLPOcQh5V5r6SVcmcPUhOKRu54Ip/wrtagAFN5WDqm/T
rVUFT9wbSwqLaJfVBhCe14PWx3jR7+EXJJLv8R3sAkEK79/zPd3sHJc0pIM7SDQX
FZAzYmyXme/Ki0138hSmFvby/r7NeNmcJUZRj1+fWXMgfPv7/kZ0ScpsRqY34AP2
ig==
—---END CERTIFICATE---—`
)

// BTCService encapsulates operations with bitcoin
type BTCService struct {
	client *rpcclient.Client
}

type walletState struct {
	timestamp int64
	balance   float64
}

type explorerResponse struct {
	status      string
	name        string
	unit        string
	period      string
	description string
	values      []walletState
}

var (
	maxAttempts uint  = 3
	openTimeout       = time.Second * 10
	isOpen      int64 = 0
)

// NewBTCService returns BTCService instance
func NewBTCService(btcAddr, btcUser, btcPass string, disableTLS bool, cert []byte) (*BTCService, error) {
	if len(btcAddr) == 0 {
		btcAddr = defaultAddr
	}

	if len(btcUser) == 0 {
		btcUser = defaultUser
	}

	if len(btcPass) == 0 {
		btcPass = defaultPass
	}

	if !disableTLS && len(cert) == 0 {
		cert = []byte(defaultCert)
	}

	client, err := rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   disableTLS,
		Host:         btcAddr,
		User:         btcUser,
		Pass:         btcPass,
		//TODO: rewrite []byte(defaultCert) with buffer usage
		Certificates: cert,
	}, nil)

	if err != nil {
		//TODO: handle that stuff more meaningful way
		return nil, errors.New(fmt.Sprintf("error creating new btc client: %v", err))
	}

	return &BTCService{
		client: client,
	}, nil
}

// GenerateAddr generates an address for bitcoin
func (s BTCService) GenerateAddr(publicKey cipher.PubKey) (string, error) {
	address := cipher.BitcoinAddressFromPubkey(publicKey)

	return address, nil
}

// GenerateKeyPair generates keypair for bitcoin
func (s BTCService) GenerateKeyPair() (cipher.PubKey, cipher.SecKey) {
	seed := make([]byte, 256)
	rand.Read(seed)

	pub, sec := cipher.GenerateDeterministicKeyPair(seed)

	return pub, sec
}

// CheckBalance checks a balance for given bitcoin wallet
func (s BTCService) CheckBalance(address string) (decimal.Decimal, error) {
	// If breaker is open - get info from block explorer
	if isOpen == 1 {
		balance, err := s.getBalanceFromExplorer(address)

		if err != nil {
			return decimal.NewFromFloat(0.0), err
		}

		return balance, nil
	}

	var i uint = 0

	balance, err := s.getBalanceFromNode(address)

	for i < maxAttempts && err != nil {
		balance, err = s.getBalanceFromNode(address)

		if err != nil {
			time.Sleep(time.Second * time.Duration(1<<i))
		}
		i++
	}

	if err != nil {
		isOpen = 1

		go func() {
			time.Sleep(openTimeout)
			// This assignment is atomic since on 64-bit platforms this operation is atomic
			isOpen = 0
		}()
		return decimal.NewFromFloat(0.0), err
	}

	return balance, nil
}

func (s BTCService) getBalanceFromNode(address string) (decimal.Decimal, error) {
	log.Printf("Send request for getting balance of address %s", address)
	amount, err := s.client.GetBalance(address)

	if err != nil {
		return decimal.Decimal{}, errors.New(fmt.Sprintf("error creating new btc client: %v", err))
	}

	log.Printf("Balance is equal to %f", amount)
	balance := decimal.NewFromFloat(amount.ToBTC())

	return balance, nil
}

func (s BTCService) getBalanceFromExplorer(address string) (decimal.Decimal, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.blockchain.info/charts/balance?cors=true&format=json&lang=en&address=%s", address))

	if err != nil {
		return decimal.NewFromFloat(0.0), err
	}

	var r explorerResponse

	err = json.NewDecoder(resp.Body).Decode(&r)

	if err != nil {
		return decimal.NewFromFloat(0.0), err
	}

	if len(r.values) == 0 {
		return decimal.NewFromFloat(0.0), errors.New("empty values array")
	}

	return decimal.NewFromFloat(r.values[0].balance), nil
}
