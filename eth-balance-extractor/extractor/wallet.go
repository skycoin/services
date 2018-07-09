package extractor

import (
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/golang-collections/go-datastructures/queue"
	"github.com/onrik/ethrpc"
)

// Wallet info
type Wallet struct {
	WalletHash        string
	PublicKey         []byte
	TransactionsCount int
	Balance           big.Int
}

// WalletScanner represents
type WalletScanner struct {
	TransactionsQueue *queue.RingBuffer
	Stop              chan int
	IsStopped         bool
	IsDisposed        bool
	Wallets           map[string]*Wallet

	TransferHash     string
	TransferFromHash string

	ProcessedTransactions []ethrpc.Transaction
	IgnoredTransactions   []ethrpc.Transaction
}

// NewWalletScanner creates a new instance of the WalletScanner
func NewWalletScanner(transactionsQueue *queue.RingBuffer, transferHash string, transferFromHash string) *WalletScanner {
	return &WalletScanner{
		TransactionsQueue:     transactionsQueue,
		Wallets:               make(map[string]*Wallet),
		IsStopped:             false,
		IsDisposed:            false,
		Stop:                  make(chan int, 10),
		TransferHash:          transferHash,
		TransferFromHash:      transferFromHash,
		ProcessedTransactions: make([]ethrpc.Transaction, 0),
		IgnoredTransactions:   make([]ethrpc.Transaction, 0),
	}
}

func hexStringToBigInt(src string) *big.Int {
	parsed := new(big.Int)
	fmt.Sscan(src, parsed)
	return parsed
}

func hexStringToBytes(src string) []byte {
	return hexStringToBigInt(src).Bytes()
}

type smartContractInput struct {
	method string

	from   string
	to     string
	amount big.Int
}

func parseSmartContractInput(input string, transferHash string, transferFromHash string) *smartContractInput {
	i := hexStringToBytes(input)
	if len(input) != 202 && len(input) != 138 {
		return nil
	}
	method := strings.ToLower(input[0:10])
	if transferHash == method {
		return &smartContractInput{
			method: strings.ToLower(input[0:10]),
			to:     strings.ToLower("0x" + input[34:74]),
			amount: *big.NewInt(0).SetBytes(i[36:68]),
		}
	}
	if transferFromHash == method {
		amount := hexStringToBigInt("0x" + input[138:202])
		return &smartContractInput{
			method: strings.ToLower(input[0:10]),
			from:   strings.ToLower("0x" + input[34:74]),
			to:     strings.ToLower("0x" + input[98:138]),
			amount: *amount,
		}
	}

	return nil
}

func normalizeHash(src []byte) []byte {
	if len(src) != 32 {
		tmp := src
		for i := 0; i < 32-len(src); i++ {
			tmp = append([]byte{0}, tmp...)
		}
		return tmp
	}
	return src
}

func recoverPublicKey(msgHash string, v string, r string, s string) ([]byte, error) {
	msgHashParsed := normalizeHash(hexStringToBytes(msgHash))
	rBytes := normalizeHash(hexStringToBytes(r))
	sBytes := normalizeHash(hexStringToBytes(s))

	vParsed := hexStringToBytes(v)[0]

	var recovery byte
	if vParsed != 27 && vParsed != 28 {
		recovery = vParsed - 35 - 2*(vParsed-35)/2
	} else {
		recovery = vParsed - 27
	}

	signature := append(append(rBytes, sBytes...), recovery)

	publicKey, err := secp256k1.RecoverPubkey(msgHashParsed, signature)
	if err != nil {
		fmt.Println()
		fmt.Println("Wallet > recoverPublicKey", err)
		fmt.Println("Sign length: ", len(signature))
		fmt.Println("Msg length: ", len(msgHashParsed), " ", msgHashParsed)
		fmt.Println("V: ", v)
		fmt.Println("R: ", len(rBytes), " ", rBytes)
		fmt.Println("S: ", len(sBytes), " ", sBytes)
		fmt.Println("================================================================================")
		return nil, err
	}

	return publicKey, nil
}

// StartScanning starts scanning process
func (w *WalletScanner) StartScanning() {
	for {
		select {
		case msg := <-w.Stop:
			if msg == 0 {
				w.IsStopped = true
			}
		default:
		}

		if w.IsStopped && w.TransactionsQueue.Len() == 0 {
			w.IsDisposed = true
			return
		}
		item, err := w.TransactionsQueue.Get()
		if item == nil {
			continue
		}
		if err != nil {
			log.Fatalln(err)
			continue
		}
		t := item.(ethrpc.Transaction)

		publicKey, _ := recoverPublicKey(t.Hash, t.V, t.R, t.S)
		input := parseSmartContractInput(t.Input, w.TransferHash, w.TransferFromHash)
		if input == nil {
			w.IgnoredTransactions = append(w.IgnoredTransactions, t)
			continue
		}

		from := ""
		if input.method == w.TransferHash {
			from = strings.ToLower(t.From)
		} else if input.method == w.TransferFromHash {
			from = strings.ToLower(input.from)
		}
		to := strings.ToLower(input.to)

		balanceFrom := big.NewInt(0)
		balanceFrom.Set(&input.amount)
		balanceFrom.Neg(balanceFrom)

		if w.Wallets[from] != nil {
			var pk []byte
			pk = nil
			if input.method == w.TransferHash {
				pk = publicKey
			}
			w.Wallets[from].Balance = *w.Wallets[from].Balance.Add(&w.Wallets[from].Balance, balanceFrom)
			w.Wallets[from].TransactionsCount++
			w.Wallets[from].PublicKey = pk
		} else {
			var pk []byte
			pk = nil
			if input.method == w.TransferHash {
				pk = publicKey
			}
			w.Wallets[from] = &Wallet{
				Balance:           *balanceFrom,
				PublicKey:         pk,
				WalletHash:        from,
				TransactionsCount: 1,
			}
		}

		balanceTo := big.NewInt(0)
		balanceTo.Set(&input.amount)

		if w.Wallets[to] != nil {
			var pk []byte
			pk = nil
			if input.method == w.TransferFromHash && to == strings.ToLower(t.From) {
				pk = publicKey
			}
			w.Wallets[to].Balance = *w.Wallets[to].Balance.Add(&w.Wallets[to].Balance, balanceTo)
			w.Wallets[to].TransactionsCount++
			w.Wallets[to].PublicKey = pk
		} else {
			var pk []byte
			pk = nil
			if input.method == w.TransferFromHash && to == strings.ToLower(t.From) {
				pk = publicKey
			}
			w.Wallets[to] = &Wallet{Balance: *balanceTo, PublicKey: pk, WalletHash: to, TransactionsCount: 1}
		}

		w.ProcessedTransactions = append(w.ProcessedTransactions, t)
	}
}

// StopScanning stops scanning process
func (w *WalletScanner) StopScanning() {
	fmt.Println("Wallet > Wallet scanner is stopping")
	w.Stop <- 0
	w.TransactionsQueue.Put(nil)
	for {
		if w.IsDisposed {
			fmt.Println("Wallet > Wallet scanner is stopped")
			return
		}
	}
}
