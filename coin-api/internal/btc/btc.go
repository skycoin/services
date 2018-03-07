package btc

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"

	"encoding/json"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/shopspring/decimal"
	"github.com/skycoin/skycoin/src/cipher"
	"io/ioutil"
	"net/http"
	"time"
	"math"
)

const (
	defaultAddr = "23.92.24.9"
	defaultUser = "YnWD3EmQAOw11IOrUJwWxAThAyobwLC"
	defaultPass = `f*Z"[1215o{qKW{Buj/wheO8@h.}j*u`
	defaultCert = `-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----`
	defaultBlockExplorer         = "https://blockchain.info"
	walletBalanceDefaultEndpoint = "/charts/balance?cors=true&format=json&lang=en&address="
	txStatusDefaultEndpoint      = "/rawtx/"
)

// ServiceBtc encapsulates operations with bitcoin
type ServiceBtc struct {
	nodeAddress string
	client      *rpcclient.Client
	// Circuit breaker related fields
	isOpen        uint32
	openTimeout   time.Duration
	retryCount    uint
	blockExplorer string
}

type walletState struct {
	Timestamp int64   `json:"x"`
	Balance   float64 `json:"y"`
}

type explorerResponse struct {
	Status      string        `json:"status"`
	Name        string        `json:"name"`
	Unit        string        `json:"unit"`
	Period      string        `json:"period"`
	Description string        `json:"description"`
	Values      []walletState `json:"values"`
}

// NewBTCService returns ServiceBtc instance
func NewBTCService(btcAddr, btcUser, btcPass string, disableTLS bool, cert []byte, blockExplorer string) (*ServiceBtc, error) {
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

	if len(blockExplorer) == 0 {
		blockExplorer = defaultBlockExplorer
	}

	client, err := rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   disableTLS,
		Host:         btcAddr,
		User:         btcUser,
		Pass:         btcPass,
		Certificates: cert,
	}, nil)

	if err != nil {
		//TODO: handle that stuff more meaningful way
		return nil, errors.New(fmt.Sprintf("error creating new btc client: %v", err))
	}

	return &ServiceBtc{
		nodeAddress:   btcAddr,
		client:        client,
		retryCount:    3,
		openTimeout:   time.Second * 10,
		isOpen:        0,
		blockExplorer: blockExplorer,
	}, nil
}

// GenerateAddr generates an address for bitcoin
func (s ServiceBtc) GenerateAddr(publicKey cipher.PubKey) (string, error) {
	address := cipher.BitcoinAddressFromPubkey(publicKey)

	return address, nil
}

// GenerateKeyPair generates keypair for bitcoin
func (s ServiceBtc) GenerateKeyPair() (cipher.PubKey, cipher.SecKey) {
	seed := make([]byte, 256)
	rand.Read(seed)

	pub, sec := cipher.GenerateDeterministicKeyPair(seed)

	return pub, sec
}

// CheckBalance checks a balance for given bitcoin wallet
func (s *ServiceBtc) CheckBalance(address string) (decimal.Decimal, error) {
	// If breaker is open - get info from block explorer
	if s.isOpen == 1 {
		balance, err := s.getBalanceFromExplorer(address)

		if err != nil {
			return decimal.NewFromFloat(0.0), err
		}

		return balance, nil
	}

	var i uint = 0

	balance, err := s.getBalanceFromNode(address)
	if err != nil {
		log.Printf("Get balance from node returned error %s", err.Error())
	}

	for i < s.retryCount && err != nil {
		if err != nil {
			log.Printf("Get balance from node returned error %s", err.Error())
		}

		balance, err = s.getBalanceFromNode(address)

		if err != nil {
			time.Sleep(time.Millisecond * time.Duration(1<<i) * 100)
		}
		i++
	}

	if err != nil {
		s.isOpen = 1

		go func() {
			time.Sleep(s.openTimeout)
			// This assignment is atomic since on 64-bit platforms this operation is atomic
			s.isOpen = 0
		}()

		balance, err := s.getBalanceFromExplorer(address)

		if err != nil {
			return decimal.NewFromFloat(0.0), err
		}

		return balance, nil
	}

	return balance, nil
}

// TODO(stgleb): Cover with unit tests
func (s *ServiceBtc) CheckTxStatus(txId string) ([]byte, error) {
	// TODO(stgleb): Add interaction with btcd node and curctui breaking on the next phase
	//hash, err := chainhash.NewHash([]byte(txId))
	//if err != nil {
	//	return nil, err
	//}
	//s.client.GetTransaction(hash)

	url := s.blockExplorer + txStatusDefaultEndpoint + txId
	resp, err := http.Get(url)

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ServiceBtc) getBalanceFromNode(address string) (decimal.Decimal, error) {
	// First get an address in proper form
	a, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)

	if err != nil {
		return decimal.NewFromFloat(0.0), err
	}

	log.Printf("Get account of address %s", address)
	account, err := s.client.GetAccount(a)

	if err != nil {
		return decimal.NewFromFloat(0.0), err
	}

	log.Printf("Send request for getting balance of address %s", address)
	amount, err := s.client.GetBalance(account)

	if err != nil {
		return decimal.Decimal{}, errors.New(fmt.Sprintf("error creating new btc client: %v", err))
	}

	log.Printf("Balance is equal to %f", amount)
	balance := decimal.NewFromFloat(amount.ToUnit(btcutil.AmountSatoshi))

	return balance, nil
}

func (s *ServiceBtc) getBalanceFromExplorer(address string) (decimal.Decimal, error) {
	url := s.blockExplorer + walletBalanceDefaultEndpoint + address
	resp, err := http.Get(url)

	if err != nil {
		return decimal.NewFromFloat(0.0), err
	}

	var r explorerResponse

	err = json.NewDecoder(resp.Body).Decode(&r)

	if err != nil {
		return decimal.NewFromFloat(0.0), err
	}

	// In case if no values - balance is empty
	if len(r.Values) == 0 {
		return decimal.NewFromFloat(0.0), nil
	}

	balance := r.Values[0].Balance / math.Pow10(int(8))
	return decimal.NewFromFloat(balance), nil
}

// Api method for monitoring btc service circuit breaker
func (s *ServiceBtc) IsOpen() bool {
	return s.isOpen == 1
}

func (s *ServiceBtc) GetHost() string {
	return s.nodeAddress
}
