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

type smartContractInput struct {
	method string
	to     string
	amount big.Int
}

func ParseSmartContractInput(input string) smartContractInput {
	i := hexStringToBytes(input)
	return smartContractInput{
		method: strings.ToLower(input[0:10]),
		to:     strings.ToLower("0x" + input[34:74]),
		amount: *big.NewInt(0).SetBytes(i[36:68]),
	}
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
			continue
		}
		t := item.(ethrpc.Transaction)

		publicKey, _ := recoverPublicKey(t.Hash, t.V, t.R, t.S)
		input := ParseSmartContractInput(t.Input)

		from := strings.ToLower(t.From)
		to := strings.ToLower(input.to)

		if input.method == "0xa9059cbb" {
			balanceFrom := big.NewInt(0)
			balanceFrom.Set(&input.amount)
			balanceFrom.Neg(balanceFrom)

			if w.Wallets[from] != nil {
				w.Wallets[from].Balance = *w.Wallets[from].Balance.Add(&w.Wallets[from].Balance, balanceFrom)
				w.Wallets[from].TransactionsCount++
			} else {
				w.Wallets[from] = &Wallet{
					Balance:           *balanceFrom,
					PublicKey:         publicKey,
					WalletHash:        from,
					TransactionsCount: 1,
				}
			}

			balanceTo := big.NewInt(0)
			balanceTo.Set(&input.amount)

			if w.Wallets[to] != nil {
				w.Wallets[to].Balance = *w.Wallets[to].Balance.Add(&w.Wallets[to].Balance, balanceTo)
				w.Wallets[to].TransactionsCount++
			} else {
				w.Wallets[to] = &Wallet{Balance: *balanceTo, PublicKey: nil, WalletHash: to, TransactionsCount: 1}
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
