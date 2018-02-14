package coin_api

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/skycoin/skycoin-exchange/_vendor-20171101171736/github.com/btcsuite/btcrpcclient"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"net/http"
)

const (
	generateKeyPair = "generateKeyPair"
	generateAddr    = "generateAddr"
	checkBalance    = "checkBalance"
	cert            = `-----BEGIN CERTIFICATE---—
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

func BtcHandler(req Request) *Response {
	switch req.Method {
	case generateKeyPair:
		return GenerateKeyPair(req)
	case generateAddr:
		return GenerateBtcAddr(req)
	case checkBalance:
		return CheckBalance(req)
	}

	return nil
}

func GenerateBtcAddr(req Request) *Response {
	if req.Params == nil {
		return &Response{
			ID:      *req.ID,
			JSONRPC: req.JSONRPC,
			Error:   MakeError(http.StatusInternalServerError, errEmptyParams.Error(), nil),
		}
	}

	result := make(map[string]interface{})
	err := json.Unmarshal(req.Params, &result)

	if err != nil {
		return &Response{
			ID:      *req.ID,
			JSONRPC: req.JSONRPC,
			Error:   MakeError(http.StatusInternalServerError, err.Error(), err),
		}
	}

	v := result["publicKey"]
	pubKey, ok := v.(cipher.PubKey)

	if !ok {
		return &Response{
			ID:      *req.ID,
			JSONRPC: req.JSONRPC,
			Error:   MakeError(http.StatusInternalServerError, err.Error(), err),
		}
	}

	addr := cipher.BitcoinAddressFromPubkey(pubKey)
	responseParams := map[string]interface{}{
		"address": addr,
	}

	return MakeSuccessResponse(req, responseParams)
}

func GenerateKeyPair(req Request) *Response {
	seed := make([]byte, 256)
	rand.Read(seed)

	pub, sec := cipher.GenerateDeterministicKeyPair(seed)

	responseParams := map[string]interface{}{
		"publicKey": pub.Hex(),
		"secretKey": sec.Hex(),
	}

	return MakeSuccessResponse(req, responseParams)
}

func CheckBalance(req Request) *Response {
	if req.Params == nil {
		return &Response{
			ID:      *req.ID,
			JSONRPC: req.JSONRPC,
			Error:   MakeError(http.StatusInternalServerError, errEmptyParams.Error(), nil),
		}
	}

	result := make(map[string]interface{})
	err := json.Unmarshal(req.Params, &result)

	if err != nil {
		return &Response{
			ID:      *req.ID,
			JSONRPC: req.JSONRPC,
			Error:   MakeError(http.StatusInternalServerError, err.Error(), err),
		}
	}

	v := result["address"]
	address, ok := v.(string)

	if !ok {
		return &Response{
			ID:      *req.ID,
			JSONRPC: req.JSONRPC,
			Error:   MakeError(http.StatusInternalServerError, err.Error(), err),
		}
	}

	balance, err := getBalance(address)

	if err != nil {
		return &Response{
			ID:      *req.ID,
			JSONRPC: req.JSONRPC,
			Error:   MakeError(http.StatusInternalServerError, err.Error(), err),
		}
	}

	responseParams := map[string]interface{}{
		"balance": balance,
	}

	return MakeSuccessResponse(req, responseParams)
}

func getBalance(address string) (decimal.Decimal, error) {
	// TODO(stgleb): Move paramas to config
	client, err := btcrpcclient.New(&btcrpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   false,
		Host:         "23.92.24.9",
		User:         `YnWD3EmQAOw11IOrUJwWxAThAyobwLC`,
		Pass:         `f*Z"[1215o{qKW{Buj/wheO8@h.}j*u`,
		Certificates: []byte(cert),
	}, nil)

	if err != nil {
		return decimal.Decimal{}, errors.New(fmt.Sprintf("error creating new btc client: %v", err))
	}

	log.Printf("Send request for getting balance of address %s", address)
	amount, err := client.GetBalance(address)
	log.Printf("Balance is equal to %f", amount)
	balance := decimal.NewFromFloat(amount.ToBTC())

	return balance, nil
}
