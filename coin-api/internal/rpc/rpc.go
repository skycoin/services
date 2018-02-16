package rpc

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
)

// JSONRPC version
const JSONRPC = "2.0"

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

func (r *Response) setBody(v interface{}) {
	body, err := json.Marshal(v)
	if err != nil {
		r.Result = nil
		r.Error = MakeError(InternalError, internalErrorMsg, err)
		return
	}
	r.Result = body
}

//Predefined messages
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)
const (
	parseErrorMsg     = "Parse Error"
	invalidRequestMsg = "Invalid Request"
	methodNotFoundMsg = "Method Not Found"
	invalidParamsMsg  = "Invalid Params"
	internalErrorMsg  = "Internal Error"
)

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

func reqID() *string {
	v, err := rand.Int(rand.Reader, new(big.Int).SetInt64(1<<62))
	if err != nil {
		panic(err)
	}
	str := v.String()
	return &str
}

func RpcRequest(addr, endpoint, method string, params map[string]interface{}) (json.RawMessage, error) {
	p, err := json.Marshal(params)
	req := Request{
		ID:      reqID(),
		JSONRPC: JSONRPC,
		Method:  method,
		Params:  p,
	}
	if err != nil {
		return nil, err
	}
	resp, err := Do(addr, endpoint, req)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

// Do does request to given addr and endpoint
func Do(addr, endpoint string, r Request) (*Response, error) {
	c := http.Client{}
	requestURI := url.URL{}

	requestURI.Host = addr
	requestURI.Scheme = "http"
	requestURI.Path = "/" + endpoint
	requestData, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, requestURI.String(), bytes.NewReader(requestData))
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			panic(err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code %d", resp.StatusCode)
	}
	respdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rpcResp Response
	err = json.Unmarshal(respdata, &rpcResp)
	if err != nil {
		return nil, err
	}
	if rpcResp.Error != nil {
		return nil, rpcResp.Error
	}
	return &rpcResp, nil

}
