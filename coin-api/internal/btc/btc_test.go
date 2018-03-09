package btc

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/skycoin/skycoin/src/cipher"
	"strings"
	"testing"
	"time"
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
		nodeAddress: "0.0.0.0",
		retryCount:  1,
		client:      client,
		isOpen:      0,
		openTimeout: time.Second * 1,
	}

	service.CheckBalance("02a1633cafcc01ebfb6d78e39f687a1f0995c62fc95f51ead10a02ee0be551b5dc")

	if !service.IsOpen() {
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
