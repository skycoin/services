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
}

// NewWalletScanner creates a new instance of the WalletScanner
func NewWalletScanner(transactionsQueue *queue.RingBuffer) *WalletScanner {
	return &WalletScanner{
		TransactionsQueue: transactionsQueue,
		Wallets:           make(map[string]*Wallet),
		IsStopped:         false,
		IsDisposed:        false,
		Stop:              make(chan int, 10),
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
		fmt.Println()
		fmt.Println()
		fmt.Println(err)
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
		}
		t := item.(ethrpc.Transaction)
		publicKey, _ := recoverPublicKey(t.Hash, t.V, t.R, t.S)

		tFrom := strings.ToLower(t.From)
		tTo := strings.ToLower(t.To)

		if tFrom != "" {
			if w.Wallets[tFrom] != nil {
				w.Wallets[tFrom].Balance = *w.Wallets[tFrom].Balance.Add(t.Value.Neg(&t.Value), &w.Wallets[tFrom].Balance)
				w.Wallets[tFrom].TransactionsCount++
			} else {
				w.Wallets[tFrom] = &Wallet{Balance: *big.NewInt(0), PublicKey: publicKey, WalletHash: t.From, TransactionsCount: 1}
			}
		}

		if tTo != "" {
			if w.Wallets[tTo] != nil {
				w.Wallets[tTo].Balance = *w.Wallets[tTo].Balance.Add(&t.Value, &w.Wallets[tTo].Balance)
				w.Wallets[tTo].TransactionsCount++
			} else {
				w.Wallets[tTo] = &Wallet{Balance: t.Value, PublicKey: nil, WalletHash: t.To, TransactionsCount: 1}
			}
		}
	}
}

// StopScanning stops scanning process
func (w *WalletScanner) StopScanning() {
	fmt.Println("Wallet scanner is stopping")
	w.Stop <- 0
	w.TransactionsQueue.Put(nil)
	for {
		if w.IsDisposed {
			fmt.Println("Wallet scanner is stopped")
			return
		}
	}
}
