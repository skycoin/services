package servd

import (
	"crypto/rand"
	"net/http"

	"encoding/json"
	"github.com/labstack/echo"
	"github.com/skycoin/services/coin-api/internal/btc"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
)

type keyPairResponse struct {
	Public  string `json:"public"`
	Private []byte `json:"private"`
}

type balanceResponse struct {
	Balance float64 `json:"balance"`
	Address string  `json:"address"`
}

type addressRequest struct {
	PublicKey string `json:"key"`
}

type addressResponse struct {
	Address string `json:"address"`
}

type handlerBTC struct {
	btcService *btc.ServiceBtc
	checker    BalanceChecker
}

type BtcStats struct {
	NodeStatus bool   `json:"node-status"`
	NodeHost   string `json:"node-host"`
}

// TODO(stgleb): Needs to be aligned with btcjson.GetTransactionResult
//type GetTransactionResult struct {
//	Amount          float64                       `json:"amount"`
//	Fee             float64                       `json:"fee,omitempty"`
//	Confirmations   int64                         `json:"confirmations"`
//	BlockHash       string                        `json:"blockhash"`
//	BlockIndex      int64                         `json:"blockindex"`
//	BlockTime       int64                         `json:"blocktime"`
//	TxID            string                        `json:"txid"`
//	WalletConflicts []string                      `json:"walletconflicts"`
//	Time            int64                         `json:"time"`
//	TimeReceived    int64                         `json:"timereceived"`
//	Details         []GetTransactionDetailsResult `json:"details"`
//	Hex             string                        `json:"hex"`
//}

type TxStatus struct {
	Ver    int `json:"ver"`
	Inputs []struct {
		Sequence int64  `json:"sequence"`
		Witness  string `json:"witness"`
		Script   string `json:"script"`
	} `json:"inputs"`
	Weight      int    `json:"weight"`
	BlockHeight int    `json:"block_height"`
	RelayedBy   string `json:"relayed_by"`
	Out         []struct {
		Spent   bool   `json:"spent"`
		TxIndex int    `json:"tx_index"`
		Type    int    `json:"type"`
		Addr    string `json:"addr,omitempty"`
		Value   int    `json:"value"`
		N       int    `json:"n"`
		Script  string `json:"script"`
	} `json:"out"`
	LockTime    int    `json:"lock_time"`
	Size        int    `json:"size"`
	DoubleSpend bool   `json:"double_spend"`
	Time        int    `json:"time"`
	TxIndex     int    `json:"tx_index"`
	VinSz       int    `json:"vin_sz"`
	Hash        string `json:"hash"`
	VoutSz      int    `json:"vout_sz"`
}

func newHandlerBTC(btcAddr, btcUser, btcPass string, disableTLS bool, cert []byte, blockExplorer string) (*handlerBTC, error) {
	log.Printf("Start new BTC handler with host %s user %s", btcAddr, btcUser)
	service, err := btc.NewBTCService(btcAddr, btcUser, btcPass, disableTLS, cert, blockExplorer)

	if err != nil {
		return nil, err
	}

	return &handlerBTC{
		btcService: service,
		checker:    service,
	}, nil
}

func (h *handlerBTC) generateKeyPair(ctx echo.Context) error {
	buffer := make([]byte, 256)
	_, err := rand.Read(buffer)

	if err != nil {
		return err
	}

	public, private := btc.ServiceBtc{}.GenerateKeyPair()

	if err := public.Verify(); err != nil {
		return handleError(ctx, err)
	}

	resp := struct {
		Status string          `json:"status"`
		Code   int             `json:"code"`
		Result keyPairResponse `json:"result"`
	}{
		"Ok",
		http.StatusOK,
		keyPairResponse{
			Public:  public.Hex(),
			Private: private[:],
		},
	}

	// Write response with newly created key pair
	ctx.JSON(http.StatusCreated, resp)
	return nil
}

func (h *handlerBTC) generateAddress(ctx echo.Context) error {
	var req addressRequest

	if err := ctx.Bind(&req); err != nil {
		return handleError(ctx, err)
	}

	if len(req.PublicKey) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "public key is empty")
	}

	publicKey, err := cipher.PubKeyFromHex(req.PublicKey)

	if err != nil {
		return handleError(ctx, err)
	}

	address, err := btc.ServiceBtc{}.GenerateAddr(publicKey)

	if err != nil {
		return handleError(ctx, err)
	}

	resp := struct {
		Status string          `json:"status"`
		Code   int             `json:"code"`
		Result addressResponse `json:"result"`
	}{
		Status: "",
		Code:   http.StatusOK,
		Result: addressResponse{
			Address: address,
		},
	}

	ctx.JSON(http.StatusCreated, resp)
	return nil
}

func (h *handlerBTC) checkTransaction(ctx echo.Context) error {
	txId := ctx.ParamValues()[0]
	data, err := h.btcService.CheckTxStatus(txId)

	if err != nil {
		return handleError(ctx, err)
	}

	status := &TxStatus{}
	err = json.Unmarshal(data, status)

	if err != nil {
		handleError(ctx, err)
		return nil
	}

	ctx.JSONPretty(http.StatusOK, struct {
		Status string   `json:"status"`
		Code   int      `json:"code"`
		Result TxStatus `json:"result"`
	}{
		Status: "",
		Code:   http.StatusOK,
		Result: *status,
	}, "\t")

	return nil
}

func (h *handlerBTC) checkBalance(ctx echo.Context) error {
	// TODO(stgleb): Check why address param is not passed
	address := ctx.ParamValues()[0]
	balance, err := h.checker.CheckBalance(address)

	if err != nil {
		return handleError(ctx, err)
	}

	balanceFloat, _ := balance.Float64()

	resp := struct {
		Status string          `json:"status"`
		Code   int             `json:"code"`
		Result balanceResponse `json:"result"`
	}{
		Status: "Ok",
		Code:   http.StatusOK,
		Result: balanceResponse{
			Balance: balanceFloat,
			Address: address,
		},
	}

	ctx.JSON(http.StatusOK, resp)
	return nil
}

// Hook for collecting stats
func (h handlerBTC) CollectStatuses(stats *Status) {
	stats.Lock()
	defer stats.Unlock()
	stats.Stats["btc"] = &BtcStats{
		NodeHost:   h.btcService.GetHost(),
		NodeStatus: h.btcService.IsOpen(),
	}
}

func handleError(ctx echo.Context, err error) error {
	return ctx.JSONPretty(http.StatusOK, struct {
		Status string `json:"status"`
		Code   int    `json:"code"`
		Result string `json:"result"`
	}{
		Status: "",
		Code:   http.StatusOK,
		Result: err.Error(),
	}, "\t")
}
