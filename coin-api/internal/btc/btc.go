package btc

import (
	"crypto/rand"
	"errors"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"bytes"
	"github.com/skycoin/skycoin/src/cipher"
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
	defaultBlockExplorer         = "https://api.blockcypher.com"
	walletBalanceDefaultEndpoint = "/v1/btc/main/addrs/"
	txStatusDefaultEndpoint      = "/v1/btc/main/txs/"
)

// ServiceBtc encapsulates operations with bitcoin
type ServiceBtc struct {
	watcherUrl string
	httpClient *http.Client

	// Circuit breaker related fields
	balanceCircuitBreaker  *CircuitBreaker
	txStatusCircuitBreaker *CircuitBreaker

	// Block explorer url
	blockExplorer string

	// How deep we will analyze blockchain for deposits on address
	blockDepth int64
}

type TxStatus struct {
	Amount        float64 `json:"amount"`
	Confirmations int64   `json:"confirmations"`
	Fee           float64 `json:"fee"`

	BlockHash  string `json:"blockhash"`
	BlockIndex int64  `json:"block_index"`

	Hash      string `json:"hash"`
	Confirmed int64  `json:"confirmed"`
	Received  int64  `json:"received"`
}

type balanceRequest struct {
	Address  string `json:"address"`
	Currency string `json:"currency"`
}

type deposit struct {
	Amount        int `json:"amount"`
	Confirmations int `json:"confirmations"`
	Height        int `json:"height"`
}

type balanceResponse struct {
	Deposits []deposit
}

type explorerTxStatus struct {
	Total         float64 `json:"total"`
	Fees          float64 `json:"fees"`
	Confirmations int64   `json:"confirmations"`

	BlockHash  string `json:"block_hash"`
	BlockIndex int64  `json:"block_index"`

	Hash      string    `json:"hash"`
	Confirmed time.Time `json:"confirmed"`
	Received  time.Time `json:"received"`
}

type explorerAddressResponse struct {
	Address            string `json:"address"`
	TotalReceived      int    `json:"total_received"`
	TotalSent          int    `json:"total_sent"`
	Balance            int64  `json:"balance"`
	UnconfirmedBalance int64  `json:"unconfirmed_balance"`
	FinalBalance       int64  `json:"final_balance"`
	NTx                int    `json:"n_tx"`
}

// NewBTCService returns ServiceBtc instance
func NewBTCService(blockExplorer string, watcherUrl string) (*ServiceBtc, error) {
	if len(blockExplorer) == 0 {
		blockExplorer = defaultBlockExplorer
	}

	service := &ServiceBtc{
		watcherUrl: watcherUrl,
		httpClient: &http.Client{
			Timeout:   time.Second * 10,
			Transport: http.DefaultTransport,
		},
		blockExplorer: blockExplorer,
	}

	balanceCircuitBreaker := NewCircuitBreaker(service.getBalanceFromWatcher,
		service.getBalanceFromExplorer,
		time.Second*10,
		time.Second*3,
		3)

	txStatusCircuitBreaker := NewCircuitBreaker(service.getTxStatusFromNode,
		service.getTxStatusFromExplorer,
		time.Second*10,
		time.Second*3,
		3)

	service.balanceCircuitBreaker = balanceCircuitBreaker
	service.txStatusCircuitBreaker = txStatusCircuitBreaker

	return service, nil
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
func (s *ServiceBtc) CheckBalance(address string) (interface{}, error) {
	return s.balanceCircuitBreaker.Do(address)
}

func (s *ServiceBtc) CheckTxStatus(txId string) (interface{}, error) {
	return s.txStatusCircuitBreaker.Do(txId)
}

func (s *ServiceBtc) getTxStatusFromNode(txId string) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (s *ServiceBtc) getTxStatusFromExplorer(txId string) (interface{}, error) {
	url := s.blockExplorer + txStatusDefaultEndpoint + txId
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	explorerResp := &explorerTxStatus{}
	err = json.Unmarshal(data, explorerResp)

	if err != nil {
		return nil, err
	}

	txStatus := &TxStatus{
		// NOTE(stgleb): amount goes in satoshis
		Amount:        explorerResp.Total,
		Confirmations: explorerResp.Confirmations,
		Fee:           explorerResp.Fees,

		BlockHash:  explorerResp.BlockHash,
		BlockIndex: explorerResp.BlockIndex,

		Hash:      explorerResp.Hash,
		Confirmed: explorerResp.Confirmed.Unix(),
		Received:  explorerResp.Received.Unix(),
	}

	return txStatus, nil
}

func (s *ServiceBtc) getBalanceFromWatcher(address string) (interface{}, error) {
	var (
		balance float64
		buffer  bytes.Buffer
	)

	reqBody := &balanceRequest{
		Address:  address,
		Currency: "BTC",
	}

	if err := json.NewEncoder(&buffer).Encode(reqBody); err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Post(s.watcherUrl+"/outputs", "application/json", &buffer)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("watcher returned an error")
	}

	balanceResp := &balanceResponse{}

	json.NewDecoder(resp.Body).Decode(balanceResp)

	if err != nil {
		return nil, err
	}

	for _, deposit := range balanceResp.Deposits {
		balance = balance + float64(deposit.Amount)
	}

	return balance, nil
}

func (s *ServiceBtc) getBalanceFromExplorer(address string) (interface{}, error) {
	url := s.blockExplorer + walletBalanceDefaultEndpoint + address
	resp, err := http.Get(url)

	if err != nil {
		return 0, err
	}

	var r explorerAddressResponse

	err = json.NewDecoder(resp.Body).Decode(&r)

	if err != nil {
		return 0, err
	}

	return r.FinalBalance, nil
}

func (s *ServiceBtc) WatcherHost() string {
	return s.watcherUrl
}

func (s *ServiceBtc) GetStatus() string {
	if s.balanceCircuitBreaker.IsOpen() || s.txStatusCircuitBreaker.IsOpen() {
		return "down"
	} else {
		return "up"
	}
}
