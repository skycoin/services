package multi

import (
	"testing"

	"github.com/skycoin/skycoin/src/visor"

	mocklib "github.com/stretchr/testify/mock"

	"github.com/skycoin/skycoin/src/api/cli"
	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"

	"github.com/skycoin/services/coin-api/internal/multi/mock"
)

const (
	rawTxID  = "bff13a47a98402ecf2d2eee40464959ad26e0ed6047de5709ffb0c0c9fc1fca5"
	rawTxStr = "dc00000000a8558b814926ed0062cd720a572bd67367aa0d01c0769ea4800adcc89cdee524010000008756e4bde4ee1c725510a6a9a308c6a90d949de7785978599a87faba601d119f27e1be695cbb32a1e346e5dd88653a97006bf1a93c9673ac59cf7b5db7e07901000100000079216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b020000000060dfa95881cdc827b45a6d49b11dbc152ecd4de640420f00000000000000000000000000006409744bcacb181bf98b1f02a11e112d7e4fa9f940f1f23a000000000000000000000000"
)

var integration bool = false

// rpcApiMck - its a mock of client web rpc API but be careful if you want laucnh tests in parallel, you may get race on
// the package-level variables and in this case you'd better off moving this variable somewhere to the getTestedService() function
var rpcApiMck *mock.WebRPCAPIMock

func TestGenerateKeyPair(t *testing.T) {
	loc := Node{
		Host: "127.0.0.1",
		Port: 6420,
	}
	skyService := NewSkyService(&loc)
	keysResponse := skyService.GenerateKeyPair()
	if len(keysResponse.Private) == 0 || len(keysResponse.Public) == 0 {
		t.Fatalf("keysResponse.Private or keysResponse.Public should not be zero length")
	}

	t.Run("TestGenerateAddress", func(t *testing.T) {
		rspAdd, err := skyService.GenerateAddr(keysResponse.Private)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(rspAdd.Address) == 0 {
			t.Fatalf("address cannot be zero lenght")
		}

		t.Run("check balance", func(t *testing.T) {
			address := rspAdd.Address
			skyServiceIsolated := getTestedMockedService()
			bRsp, err := skyServiceIsolated.CheckBalance(address)
			if err != nil {
				t.Fatal(err.Error())
			}
			if bRsp.Balance != "23" || bRsp.Hours != "3" {
				t.Fatalf("wrong balance")
			}
		})
	})
}

func TestTransaction(t *testing.T) {
	skyService := getTestedMockedService()
	t.Run("sign transaction", func(t *testing.T) {
		_, secKey := makeUxBodyWithSecret(t)
		secKeyHex := secKey.Hex()
		bRsp, err := skyService.SignTransaction(secKeyHex, rawTxStr)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(bRsp.Signid) == 0 {
			t.Fatalf("signid shouldn't be zero length")
		}
	})

	t.Run("inject transaction", func(t *testing.T) {
		// testing doubles: test input and generate output
		rpcApiMck.On("GetTransactionByID", mocklib.MatchedBy(func(txid string) bool {
			if rawTxID != txid {
				return false
			}
			return true
		})).Return(&webrpc.TxnResult{
			Transaction: &visor.TransactionResult{
				Status: visor.TransactionStatus{
					Confirmed: true,
					Height:    12799,
					BlockSeq:  12799,
				},
				Time: 127993444,
			},
		}, nil)

		rpcApiMck.On("InjectTransactionString", mocklib.MatchedBy(func(txid string) bool {
			if rawTxStr != txid {
				return false
			}
			return true
		})).Return(rawTxID, nil)

		bRsp, err := skyService.InjectTransaction(rawTxStr)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(bRsp.Transid) == 0 {
			t.Fatalf("signid shouldn't be zero length")
		}
	})

	t.Run("check transaction status", func(t *testing.T) {
		rpcApiMck.On("GetTransactionByID", mocklib.MatchedBy(func(txid string) bool {
			if rawTxID != txid {
				return false
			}
			return true
		})).Return(&webrpc.TxnResult{
			Transaction: &visor.TransactionResult{
				Status: visor.TransactionStatus{
					Confirmed: true,
					Height:    12799,
					BlockSeq:  12799,
				},
				Time: 127993444,
			},
		}, nil)

		transStatus, err := skyService.CheckTransactionStatus(rawTxID)

		if err != nil {
			t.Fatal(err.Error())
		}

		if transStatus.BlockSeq == 0 {
			t.Fatalf("blockSeq shouldn't be zero length")
		}
	})
}

var getTestedMockedService = func() *SkyСoinService {
	loc := Node{
		Host: "127.0.0.1",
		Port: 6430,
	}
	// parametrize tested service with mocked/stubbed external services
	// this way we mock our helpers which commit 3-d party package calls which cannot be mocked usual way because they deal with types
	// instead of interfaces
	getBalanceAddresses := func(client WebRPCClientAPI, addresses []string) (*cli.BalanceResult, error) {
		return &cli.BalanceResult{
			Confirmed: cli.Balance{
				Coins: "23",
				Hours: "3",
			},
			Spendable: cli.Balance{
				Coins: "23",
				Hours: "3",
			},
		}, nil
	}

	rpcApiMck = &mock.WebRPCAPIMock{}
	skyService := NewSkyService(&loc)
	// inject mocked dependencies into tested service
	skyService.InjectRPCAPIMock(rpcApiMck)
	skyService.InjectCheckBalanceMock(getBalanceAddresses)

	return skyService
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
