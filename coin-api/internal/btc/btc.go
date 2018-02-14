package btc

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
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

// Request represents a JSONRPC 2.0 request message
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      *string         `json:"id"`
}

// Response represents a JSONRPC 2.0 response message
type Response struct {
	ID      string          `json:"id,omitempty"`
	Error   *jsonrpcError   `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	JSONRPC string          `json:"jsonrpc"`
}

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

// MakeSuccessResponse creates success response
func MakeSuccessResponse(r Request, result interface{}) *Response {
	data, err := json.Marshal(result)
	if err != nil {
		return &Response{
			ID:      *r.ID,
			JSONRPC: r.JSONRPC,
			Error:   MakeError(InternalError, internalErrorMsg, err),
		}
	}
	return &Response{
		ID:      *r.ID,
		JSONRPC: JSONRPC,
		Result:  data,
	}
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

type jsonrpcError struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    *string `json:"data,omitempty"`
}

// Implements error interface
func (err *jsonrpcError) Error() string {
	return fmt.Sprintf("jsonrpc error: %d %s %s", err.Code, err.Message, *err.Data)
}

func MakeError(code int, message string, additional error) *jsonrpcError {
	var datastr *string
	if additional != nil {
		datastr = new(string)
		*datastr = additional.Error()
	}
	return &jsonrpcError{Code: code, Message: message, Data: datastr}
}
