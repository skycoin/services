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

// TransactionWallet represents wallet info with transaction that will be used to restore public key
type TransactionWallet struct {
	WalletHash        string
	TxHash            string
	TransactionsCount int
	Balance           big.Int
}

// Wallet info
type Wallet struct {
	WalletHash        string
	PublicKey         []byte
	TransactionsCount int
	Balance           big.Int
}

// WalletScanner represents
type WalletScanner struct {
	nodeAPIURL         string
	Wallets            map[string]*Wallet
	TransactionWallets map[string]*TransactionWallet
	TransactionsQueue  *queue.RingBuffer

	Stop       chan int
	IsStopped  bool
	IsDisposed bool
}

// NewWalletScanner creates a new instance of the WalletScanner
func NewWalletScanner(nodeAPIURL string, transactionsQueue *queue.RingBuffer, transactionWallets map[string]*TransactionWallet) *WalletScanner {
	return &WalletScanner{
		nodeAPIURL:         nodeAPIURL,
		Wallets:            make(map[string]*Wallet),
		TransactionWallets: transactionWallets,
		TransactionsQueue:  transactionsQueue,

		IsStopped:  false,
		IsDisposed: false,
		Stop:       make(chan int, 10),
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

// RestoreKeys starts process of restoring public keys
func (w *WalletScanner) RestoreKeys() {
	client := ethrpc.New(w.nodeAPIURL)

	iteration := 0
	total := len(w.TransactionWallets)

	for key, tw := range w.TransactionWallets {
		iteration = iteration + 1
		fmt.Println("Iteration", iteration, "of", total)

		if tw.TxHash == "" {
			w.Wallets[key] = &Wallet{
				Balance:           tw.Balance,
				PublicKey:         nil,
				TransactionsCount: tw.TransactionsCount,
				WalletHash:        key,
			}
		} else {
			tx, err := client.EthGetTransactionByHash(tw.TxHash)
			if err != nil {
				w.Wallets[key] = &Wallet{
					Balance:           tw.Balance,
					PublicKey:         nil,
					TransactionsCount: tw.TransactionsCount,
					WalletHash:        key,
				}
				fmt.Println("WalletScanner > RestoreKeys", err)
			} else {
				if tx.From != tw.WalletHash {
					w.Wallets[key] = &Wallet{
						Balance:           tw.Balance,
						PublicKey:         nil,
						TransactionsCount: tw.TransactionsCount,
						WalletHash:        key,
					}
					fmt.Println("Wrong transaction:", tx.From, tw.WalletHash)
					continue
				}

				pk, err := recoverPublicKey(tx.Hash, tx.V, tx.R, tx.S)
				if err != nil {
					w.Wallets[key] = &Wallet{
						Balance:           tw.Balance,
						PublicKey:         nil,
						TransactionsCount: tw.TransactionsCount,
						WalletHash:        key,
					}
					fmt.Println("WalletScanner > RestoreKeys", err)
				}

				w.Wallets[key] = &Wallet{
					Balance:           tw.Balance,
					PublicKey:         pk,
					TransactionsCount: tw.TransactionsCount,
					WalletHash:        key,
				}
			}
		}
	}
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

		from := strings.ToLower(t.From)

		if w.Wallets[from] == nil {
			continue
		}

		publicKey, _ := recoverPublicKey(t.Hash, t.V, t.R, t.S)
		w.Wallets[from].PublicKey = publicKey
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
