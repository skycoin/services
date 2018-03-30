package btc

import (
	"strings"
	"testing"
	"time"

	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
)

func TestCheckBalanceOpen(t *testing.T) {
	service := ServiceBtc{
		watcherUrl:    "http://localhost:8080",
		httpClient:    http.DefaultClient,
		blockExplorer: "https://api.blockcypher.com",
	}

	// Create circuit breaker for btc service
	balanceCircuitBreaker := NewCircuitBreaker(service.getBalanceFromWatcher,
		service.getBalanceFromExplorer,
		time.Second*10,
		time.Second*3,
		3)
	service.balanceCircuitBreaker = balanceCircuitBreaker

	service.CheckBalance("02a1633cafcc01ebfb6d78e39f687a1f0995c62fc95f51ead10a02ee0be551b5dc")

	if !service.balanceCircuitBreaker.IsOpen() {
		t.Error("Expected curcuit breaker to be open, actual closed")
	}
}

func TestCheckBalanceClosed(t *testing.T) {
	service := ServiceBtc{
		watcherUrl:    "http://localhost:8080",
		httpClient:    http.DefaultClient,
		blockExplorer: "https://api.blockcypher.com",
	}
	expectedBalance := int64(5)
	expectedAddress := "02a1633cafcc01ebfb6d78e39f687a1f0995c62fc95f51ead10a02ee0be551b5dc"
	expectedDepositCount := 2

	success := func(add string) (interface{}, error) {
		return &BalanceResponse{
			Address: expectedAddress,
			Balance: expectedBalance,
			Utxo: []Deposit{
				{
					Amount:        2,
					Confirmations: 8,
					Height:        9,
				},
				{
					Amount:        3,
					Confirmations: 7,
					Height:        10,
				},
			},
		}, nil
	}

	fallback := func(add string) (interface{}, error) {
		return nil, nil
	}

	// Create circuit breaker for btc service
	balanceCircuitBreaker := NewCircuitBreaker(
		success,
		fallback,
		time.Second*10,
		time.Second*3,
		3)
	service.balanceCircuitBreaker = balanceCircuitBreaker
	resp, err := service.CheckBalance(expectedAddress)

	if service.balanceCircuitBreaker.IsOpen() {
		t.Error("Expected curcuit breaker to be closed, actual open")
	}

	if err != nil {
		t.Error(err)
	}

	balance, ok := resp.(*BalanceResponse)

	if !ok {
		t.Errorf("Wrong type conversion expected *BalanceResponse actual %T", resp)
	}

	if expectedBalance != balance.Balance {
		t.Errorf("Wrong balance  expected %d actual %d", expectedBalance, balance)
	}

	if expectedAddress != balance.Address {
		t.Errorf("Wrong address expected %s actual %s", expectedAddress, balance.Address)
	}

	if expectedDepositCount != len(balance.Utxo) {
		t.Errorf("Wrong Deposit count expected %d actual %d", expectedDepositCount, len(balance.Utxo))
	}
}

func TestServiceBtcCheckTxStatus(t *testing.T) {
	service := ServiceBtc{
		watcherUrl:    "http://localhost:8080",
		httpClient:    http.DefaultClient,
		blockExplorer: "https://api.blockcypher.com",
	}

	txStatusCircuitBreaker := NewCircuitBreaker(service.getTxStatusFromNode,
		service.getTxStatusFromExplorer,
		time.Second*10,
		time.Second*3,
		3)
	service.txStatusCircuitBreaker = txStatusCircuitBreaker

	result, err := service.CheckTxStatus("f854aebae95150b379cc1187d848d58225f3c4157fe992bcd166f58bd5063449")

	if err != nil {
		t.Error(err)
		return
	}

	status, _ := result.(*TxStatus)

	if status.Confirmations == 0 {
		t.Errorf("txRef confirmations must be greater than zero")
		return
	}

	if !service.txStatusCircuitBreaker.IsOpen() {
		t.Error("Expected curcuit breaker to be open, actual closed")
	}
}

func TestGenerateAddr(t *testing.T) {
	service := &ServiceBtc{}
	key := "02a1633cafcc01ebfb6d78e39f687a1f0995c62fc95f51ead10a02ee0be551b5dc"
	pk, err := cipher.PubKeyFromHex(key)

	if err != nil {
		t.Error(err)
	}

	addr, err := service.GenerateAddr(pk)

	if err != nil {
		t.Error(err)
	}

	if strings.Compare(addr, "17JarKo61PkpuZG3GyofzGmFSCskGRBUT3") != 0 {
		t.Error("wrong address value")
	}
}
