package coin_api

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/shopspring/decimal"
	"github.com/skycoin/services/coin-api/internal/rpc"
	"github.com/skycoin/skycoin/src/cipher"
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

var (
	errEmptyParams = errors.New("empty params")
	errWrongType   = errors.New("wrong type")
)

// BTCService encapsulates operations with bitcoin
type BTCService struct {
	client *rpcclient.Client
}

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
func (s *BTCService) GenerateAddr(req rpc.Request) *rpc.Response {
	if req.Params == nil {
		return &rpc.Response{
			ID:      *req.ID,
			JSONRPC: req.JSONRPC,
			Error:   rpc.MakeError(http.StatusInternalServerError, errEmptyParams.Error(), nil),
		}
	}

	result := make(map[string]interface{})
	err := json.Unmarshal(req.Params, &result)

	if err != nil {
		return &rpc.Response{
			ID:      *req.ID,
			JSONRPC: req.JSONRPC,
			Error:   rpc.MakeError(http.StatusInternalServerError, err.Error(), err),
		}
	}

	v := result["publicKey"]
	pubKey, ok := v.(cipher.PubKey)

	if !ok {
		return &rpc.Response{
			ID:      *req.ID,
			JSONRPC: req.JSONRPC,
			Error:   rpc.MakeError(http.StatusInternalServerError, err.Error(), err),
		}
	}

	addr := cipher.BitcoinAddressFromPubkey(pubKey)
	responseParams := map[string]interface{}{
		"address": addr,
	}

	return rpc.MakeSuccessResponse(req, responseParams)
}

// GenerateKeyPair generates keypair for bitcoin
func (s *BTCService) GenerateKeyPair(req rpc.Request) *rpc.Response {
	seed := make([]byte, 256)
	rand.Read(seed)

	pub, sec := cipher.GenerateDeterministicKeyPair(seed)

	responseParams := map[string]interface{}{
		"publicKey": pub.Hex(),
		"secretKey": sec.Hex(),
	}

	return rpc.MakeSuccessResponse(req, responseParams)
}

// CheckBalance checks a balance for given bitcoin wallet
func (s *BTCService) CheckBalance(req rpc.Request) *rpc.Response {
	if req.Params == nil {
		return &rpc.Response{
			ID:      *req.ID,
			JSONRPC: req.JSONRPC,
			Error:   rpc.MakeError(http.StatusInternalServerError, errEmptyParams.Error(), nil),
		}
	}

	result := make(map[string]interface{})
	err := json.Unmarshal(req.Params, &result)

	if err != nil {
		return &rpc.Response{
			ID:      *req.ID,
			JSONRPC: req.JSONRPC,
			Error:   rpc.MakeError(http.StatusInternalServerError, err.Error(), err),
		}
	}

	v := result["address"]
	address, ok := v.(string)

	if !ok {
		return &rpc.Response{
			ID:      *req.ID,
			JSONRPC: req.JSONRPC,
			Error:   rpc.MakeError(http.StatusInternalServerError, err.Error(), err),
		}
	}

	balance, err := s.getBalance(address)

	if err != nil {
		return &rpc.Response{
			ID:      *req.ID,
			JSONRPC: req.JSONRPC,
			Error:   rpc.MakeError(http.StatusInternalServerError, err.Error(), err),
		}
	}

	responseParams := map[string]interface{}{
		"balance": balance,
	}

	return rpc.MakeSuccessResponse(req, responseParams)
}

func (s *BTCService) getBalance(address string) (decimal.Decimal, error) {
	log.Printf("Send request for getting balance of address %s", address)
	amount, err := s.client.GetBalance(address)

	if err != nil {
		return decimal.Decimal{}, errors.New(fmt.Sprintf("error creating new btc client: %v", err))
	}

	log.Printf("Balance is equal to %f", amount)
	balance := decimal.NewFromFloat(amount.ToBTC())

	return balance, nil
}
