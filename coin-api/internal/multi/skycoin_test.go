package multi_test

import (
	"testing"

	"reflect"

	"github.com/stretchr/testify/assert"

	"github.com/skycoin/services/coin-api/internal/locator"
	"github.com/skycoin/services/coin-api/internal/model"
	"github.com/skycoin/services/coin-api/internal/multi"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
)

const (
	// rawTxID = "bdc4a85a3e9d17a8fe00aa7430d0347c7f1dd6480a16da7147b6e43905057d43"
	rawTxID  = "bff13a47a98402ecf2d2eee40464959ad26e0ed6047de5709ffb0c0c9fc1fca5"
	rawTxStr = "dc00000000a8558b814926ed0062cd720a572bd67367aa0d01c0769ea4800adcc89cdee524010000008756e4bde4ee1c725510a6a9a308c6a90d949de7785978599a87faba601d119f27e1be695cbb32a1e346e5dd88653a97006bf1a93c9673ac59cf7b5db7e07901000100000079216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b020000000060dfa95881cdc827b45a6d49b11dbc152ecd4de640420f00000000000000000000000000006409744bcacb181bf98b1f02a11e112d7e4fa9f940f1f23a000000000000000000000000"
)

func TestGenerateAddress(t *testing.T) {
	loc := locator.Node{
		Host: "127.0.0.1",
		Port: 6430,
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
		address := rspAdd.Address
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
			t.Fatalf("Address shouldn't be zero length")
		}
	})
}

func TestTransaction(t *testing.T) {
	loc := locator.Node{
		Host: "127.0.0.1",
		// Port: 20200,
		Port: 6430,
	}

	skyService := multi.NewSkyService(&loc)
	t.Run("sign transaction", func(t *testing.T) {
		//TODO: check this logic
		_, secKey := makeUxBodyWithSecret(t)
		secKeyHex := secKey.Hex()
		rsp, err := skyService.SignTransaction(secKeyHex, rawTxStr)
		if !assert.NoError(t, err) {
			println("err.Error()", err.Error())
			t.FailNow()
		}
		assertCodeZero(t, rsp)
		assertStatusOk(t, rsp)
		result := rsp.Result
		bRsp, ok := result.(*model.TransactionSign)
		if !ok {
			t.Fatalf("wrong type, *model.TransactionSign expected, given %s", reflect.TypeOf(result).String())
		}
		if len(bRsp.Signid) == 0 {
			t.Fatalf("signid shouldn't be zero length")
		}
	})

	t.Run("inject transaction", func(t *testing.T) {
		rsp, err := skyService.InjectTransaction(rawTxStr)
		if !assert.NoError(t, err) {
			println("err.Error()", err.Error())
			t.FailNow()
		}
		assertCodeZero(t, rsp)
		assertStatusOk(t, rsp)
		result := rsp.Result
		bRsp, ok := result.(*model.Transaction)
		if !ok {
			t.Fatalf("wrong type, *model.Transaction expected, given %s", reflect.TypeOf(result).String())
		}
		if len(bRsp.Transid) == 0 {
			t.Fatalf("signid shouldn't be zero length")
		}
	})

	t.Run("check transaction status", func(t *testing.T) {
		transStatus, err := skyService.CheckTransactionStatus(rawTxID)
		if !assert.NoError(t, err) {
			println("err.Error()", err.Error())
			t.FailNow()
		}
		if transStatus.BlockSeq == 0 {
			t.Fatalf("blockSeq shouldn't be zero length")
		}
	})
}

func TestGenerateKeyPair(t *testing.T) {
	loc := locator.Node{
		Host: "127.0.0.1",
		Port: 6420,
	}
	skyService := multi.NewSkyService(&loc)
	rsp := skyService.GenerateKeyPair()
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

func makeUxBodyWithSecret(t *testing.T) (coin.UxBody, cipher.SecKey) {
	p, s := cipher.GenerateKeyPair()
	return coin.UxBody{
		SrcTransaction: testutil.RandSHA256(t),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          1e6,
		Hours:          100,
	}, s
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
