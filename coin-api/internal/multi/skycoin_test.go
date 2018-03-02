package multi_test

import (
	"testing"

	"reflect"

	"github.com/stretchr/testify/assert"

	"github.com/skycoin/services/coin-api/internal/locator"
	"github.com/skycoin/services/coin-api/internal/model"
	"github.com/skycoin/services/coin-api/internal/multi"
)

func TestGenerateAddress(t *testing.T) {
	loc := locator.Node{
		Host: "127.0.0.1",
		// Port: 6420,
		Port: 6000,
	}
	skyService := multi.NewSkyService(&loc)
	rsp, err := skyService.GenerateAddr(1, true)
	assert.NoError(t, err)
	assertCodeZero(t, rsp)
	assertStatusOk(t, rsp)
	result := rsp.Result
	rspAdd, ok := result.(*model.AddressResponse)
	if !ok {
		t.Fatalf("wrong type, result.(*model.AddressResponse) expected, given %s", reflect.TypeOf(result).String())
	}
	if len(rspAdd.Address) == 0 {
		t.Fatalf("address cannot be zero lenght")
	}
	t.Run("check balance", func(t *testing.T) {
		// walletFile := "somefile.wlt"
		// addr := 4029003020
		address := rspAdd.Address
		address = "LinfYSSC8cK13mA3KYJ2xczEsdFABbB2yP"
		rsp, err := skyService.CheckBalance(address)
		if !assert.NoError(t, err) {
			t.Fatal()
		}
		assertCodeZero(t, rsp)
		assertStatusOk(t, rsp)
		result := rsp.Result
		bRsp, ok := result.(*model.BalanceResponse)
		if !ok {
			t.Fatalf("wrong type, *model.BalanceResponse expected, given %s", reflect.TypeOf(result).String())
		}
		if len(bRsp.Address) == 0 {
			t.Fatalf("bRsp.Address shouldn't be zero ")
		}
	})
}

func TestGenerateKeyPair(t *testing.T) {
	loc := locator.Node{
		Host: "127.0.0.1",
		Port: 6420,
	}
	sky := multi.NewSkyService(&loc)
	rsp := sky.GenerateKeyPair()
	assertCodeZero(t, rsp)
	assertStatusOk(t, rsp)
	result := rsp.Result
	keysResponse, ok := result.(*model.KeysResponse)
	if !ok {
		t.Fatalf("wrong type, result.(*model.KeysResponse) expected, given %s", reflect.TypeOf(result).String())
	}
	if len(keysResponse.Private) == 0 || len(keysResponse.Public) == 0 {
		t.Fatalf("keysResponse.Private or keysResponse.Public should not be zero length")
	}
}

func assertCodeZero(t *testing.T, rsp *model.Response) {
	if rsp.Code != 0 {
		t.Fatalf("the code must be 0, %d given", rsp.Code)
	}
}

func assertStatusOk(t *testing.T, rsp *model.Response) {
	if rsp.Status != "ok" {
		t.Fatalf("status must be ok %s given", rsp.Status)
	}
}
