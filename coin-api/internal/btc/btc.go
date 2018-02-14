package coin_api

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
	"net/http"
)

const (
	generateKeyPair = "generateKeyPair"
	generateAddr    = "generateAddr"
	checkBalance    = "checkBalance"
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

	// TODO(stgleb): Find api how btc addr balance can be checked
	balance := len(address)
	responseParams := map[string]interface{}{
		"balance": balance,
	}

	return MakeSuccessResponse(req, responseParams)
}
