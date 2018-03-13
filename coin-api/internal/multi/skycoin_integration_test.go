package multi_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"

// 	"github.com/skycoin/services/coin-api/internal/locator"
// 	"github.com/skycoin/services/coin-api/internal/multi"
// )

// func TestTransactionIntegration(t *testing.T) {
// 	skyService := getTestedService()
// 	t.Run("sign transaction", func(t *testing.T) {
// 		_, secKey := makeUxBodyWithSecret(t)
// 		secKeyHex := secKey.Hex()
// 		bRsp, err := skyService.SignTransaction(secKeyHex, rawTxStr)
// 		if !assert.NoError(t, err) {
// 			t.Fatal(err.Error())
// 		}
// 		if len(bRsp.Signid) == 0 {
// 			t.Fatalf("signid shouldn't be zero length")
// 		}
// 	})

// 	t.Run("inject transaction", func(t *testing.T) {
// 		bRsp, err := skyService.InjectTransaction(rawTxStr)
// 		if !assert.NoError(t, err) {
// 			t.Fatal(err.Error())
// 		}
// 		if len(bRsp.Transid) == 0 {
// 			t.Fatalf("signid shouldn't be zero length")
// 		}
// 	})

// 	t.Run("check transaction status", func(t *testing.T) {
// 		transStatus, err := skyService.CheckTransactionStatus(rawTxID)
// 		if !assert.NoError(t, err) {
// 			t.Fatal(err.Error())
// 		}
// 		if transStatus.BlockSeq == 0 {
// 			t.Fatalf("blockSeq shouldn't be zero length")
// 		}
// 	})
// }

// var getTestedService = func() *multi.SkyСoinService {
// 	loc := locator.Node{
// 		Host: "127.0.0.1",
// 		Port: 6430,
// 	}

// 	return multi.NewSkyService(&loc)
// }
