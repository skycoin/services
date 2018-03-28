package btc

import (
	"crypto/rand"
	"errors"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"bytes"

	"fmt"
	"github.com/skycoin/skycoin/src/cipher"
)

const (
	defaultBlockExplorer         = "https://api.blockcypher.com"
	walletBalanceDefaultEndpoint = "/v1/btc/main/addrs/"
	txDefaultEndpoint            = "/v1/btc/main/txs/"
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
		buffer bytes.Buffer
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
		return nil, fmt.Errorf("watcher returned an error %d", resp.StatusCode)
	}

	balanceResp := &BalanceResponse{}

	json.NewDecoder(resp.Body).Decode(balanceResp)

	if err != nil {
		return nil, err
	}

	// Summarize Deposit values
	for _, deposit := range balanceResp.Utxo {
		balanceResp.Balance += int64(deposit.Amount)
	}

	balanceResp.Address = address

	return balanceResp, nil
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

	balanceResp := &BalanceResponse{
		Address:     address,
		Balance:     r.FinalBalance,
		Utxo:        make([]Deposit, 0),
		PendingUtxo: make([]Deposit, 0),
	}

	// Collect input transactions for the address,
	// see for detail https://blockcypher.github.io/documentation/#transactions
	for _, tx := range r.Transactions {
		if tx.TxInputN == -1 {
			blockHash := getTxBlockHash(s.blockExplorer, tx.TxHash)

			dep := Deposit{
				tx.Value,
				tx.TxHash,
				blockHash,
				tx.Confirmations,
				tx.BlockHeight,
			}
			balanceResp.Utxo = append(balanceResp.Utxo, dep)
		}
	}

	// Collect pending incoming transactions
	for _, tx := range r.UnconfirmedTransactions {
		if tx.TxInputN == -1 {
			blockHash := getTxBlockHash(s.blockExplorer, tx.TxHash)

			dep := Deposit{
				tx.Value,
				tx.TxHash,
				blockHash,
				tx.Confirmations,
				tx.BlockHeight,
			}
			balanceResp.PendingUtxo = append(balanceResp.PendingUtxo, dep)
		}
	}

	return balanceResp, nil
}

func getTxBlockHash(blockExplorer, txHash string) string {
	var txInfo TxInfo

	txUrl := blockExplorer + txDefaultEndpoint + txHash
	resp, _ := http.Get(txUrl)

	json.NewDecoder(resp.Body).Decode(&txInfo)

	return txInfo.BlockHash
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
