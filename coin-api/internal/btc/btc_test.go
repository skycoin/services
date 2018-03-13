package btc

import (
	"strings"
	"testing"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/skycoin/skycoin/src/cipher"
)

func TestCheckBalance(t *testing.T) {
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         "0.0.0.0",
		User:         "",
		Pass:         "",
	}, nil)

	if err != nil {
		t.Error(err)
	}

	service := ServiceBtc{
		nodeAddress:   "0.0.0.0",
		client:        client,
		blockExplorer: "https://api.blockcypher.com",
	}

	// Create circuit breaker for btc service
	balanceCircuitBreaker := NewCircuitBreaker(service.getBalanceFromNode,
		service.getBalanceFromExplorer,
		time.Second*10,
		3)
	service.balanceCircuitBreaker = balanceCircuitBreaker

	service.CheckBalance("02a1633cafcc01ebfb6d78e39f687a1f0995c62fc95f51ead10a02ee0be551b5dc")

	if !service.balanceCircuitBreaker.IsOpen() {
		t.Error("Expected curcuit breaker to be open, actual closed")
	}
}

func TestServiceBtcCheckTxStatus(t *testing.T) {
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         "0.0.0.0",
		User:         "",
		Pass:         "",
	}, nil)

	if err != nil {
		t.Error(err)
		return
	}

	service := ServiceBtc{
		nodeAddress:   "0.0.0.0",
		client:        client,
		blockExplorer: "https://api.blockcypher.com",
	}

	txStatusCircuitBreaker := NewCircuitBreaker(service.getTxStatusFromNode,
		service.getTxStatusFromExplorer,
		time.Second*10,
		3)
	service.txStatusCircuitBreaker = txStatusCircuitBreaker

	result, err := service.CheckTxStatus("f854aebae95150b379cc1187d848d58225f3c4157fe992bcd166f58bd5063449")

	if err != nil {
		t.Error(err)
		return
	}

	status, _ := result.(*TxStatus)

	if status.Confirmations == 0 {
		t.Errorf("Transaction confirmations must be greater than zero")
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
