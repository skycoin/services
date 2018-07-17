package extractor

import (
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/onrik/ethrpc"
)

// Storage represets a logic of storing wallets
type Storage struct {
	dataDir string
}

// NewStorage creates a new instance of the Storage
func NewStorage(dataDir string) *Storage {
	return &Storage{
		dataDir: dataDir,
	}
}

// StoreSnapshot saves snapshot into the CSV file
func (s *Storage) StoreSnapshot(filename string, snapshot map[string]*Wallet) {
	f, err := os.Create(fmt.Sprintf("%s/%s.csv", s.dataDir, filename))
	if err != nil {
		fmt.Println("Storage > StoreSnapshot", err)
		panic(err)
	}
	defer f.Close()
	for key, w := range snapshot {
		f.WriteString(fmt.Sprintf("%s,%s,%s,%d\n", key, hex.EncodeToString(w.PublicKey), w.Balance.Text(10), w.TransactionsCount))
	}
}

// LoadSnapshot loads a snapshot from the specified CSV file
func (s *Storage) LoadSnapshot(snapshotPath string) map[string]*Wallet {
	file, err := os.Open(snapshotPath)
	defer file.Close()
	if err != nil {
		fmt.Println("Storage > LoadSnapshot", err)
		panic(err)
	}

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Storage > LoadSnapshot", err)
		panic(err)
	}

	wallets := make(map[string]*Wallet, len(data))
	for _, d := range data {
		transactionsCount, err := strconv.Atoi(d[3])
		if err != nil {
			fmt.Println("Storage > LoadSnapshot", err)
			panic(err)
		}

		var pubKey []byte
		pubKey = nil
		if d[1] != "" {
			pubKey, err = hex.DecodeString(d[1])
			if err != nil {
				fmt.Println("Storage > LoadSnapshot", err)
				panic(err)
			}
		}

		balance := big.NewInt(0)
		balance.SetString(d[2], 10)

		w := &Wallet{
			WalletHash:        d[0],
			PublicKey:         pubKey,
			Balance:           *balance,
			TransactionsCount: transactionsCount,
		}
		wallets[d[0]] = w
	}

	return wallets
}

// LoadTransactionWallets loads a snapshot from the specified CSV file
func (s *Storage) LoadTransactionWallets(snapshotPath string) map[string]*TransactionWallet {
	file, err := os.Open(snapshotPath)
	defer file.Close()
	if err != nil {
		fmt.Println("Storage > LoadSnapshot", err)
		panic(err)
	}

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Storage > LoadSnapshot", err)
		panic(err)
	}

	wallets := make(map[string]*TransactionWallet, len(data))
	for _, d := range data {
		transactionsCount, err := strconv.Atoi(d[3])
		if err != nil {
			fmt.Println("Storage > LoadSnapshot", err)
			panic(err)
		}

		var txHash string
		if d[1] == "undefined" || d[1] == "null" {
			txHash = ""
		} else {
			txHash = d[1]
		}

		balance := big.NewInt(0)
		balance.SetString(d[2], 10)

		w := &TransactionWallet{
			WalletHash:        d[0],
			TxHash:            txHash,
			Balance:           *balance,
			TransactionsCount: transactionsCount,
		}
		wallets[d[0]] = w
	}

	return wallets
}

// SaveTransactions safes an array of transactions
func (s *Storage) SaveTransactions(blocksScanned int, folderName string, transactions []ethrpc.Transaction) {
	f, err := os.Create(fmt.Sprintf("%s/%s/%d.json", s.dataDir, folderName, blocksScanned))
	if err != nil {
		fmt.Println("Storage > StoreSnapshot", err)
		panic(err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(transactions)
}
