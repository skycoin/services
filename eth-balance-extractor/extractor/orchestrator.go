package extractor

import (
	"fmt"
)

// Orchestrator represents an orchestrator of transactions extractor and wallets scanner
type Orchestrator struct {
	extractor        *Extractor
	scanner          *WalletScanner
	storage          *Storage
	contractAddress  string
	transferHash     string
	transferFromHash string
	nodeAPI          string

	startBlock   int
	threadsCount int
}

// NewOrchestrator creates a new instance of the Orchestrator
func NewOrchestrator(nodeAPI string, contractAddress string, transferHash string, transferFromHash string, destDir string, startBlock int, threadsCount int) *Orchestrator {
	e := NewExtractor(nodeAPI, contractAddress)
	return &Orchestrator{
		extractor:        e,
		scanner:          NewWalletScanner(e.TransactionsQueue, transferHash, transferFromHash),
		storage:          NewStorage(destDir),
		contractAddress:  contractAddress,
		transferHash:     transferHash,
		transferFromHash: transferFromHash,
		nodeAPI:          nodeAPI,

		startBlock:   startBlock,
		threadsCount: threadsCount,
	}
}

func reinitializeOrchestrator(o *Orchestrator, wallets map[string]*Wallet) {
	e := NewExtractor(o.nodeAPI, o.contractAddress)
	o.extractor = e
	o.scanner = NewWalletScanner(e.TransactionsQueue, o.transferHash, o.transferFromHash)
	o.scanner.Wallets = wallets
}

// StartScanning starts scanning process
func (o *Orchestrator) StartScanning() {
	lastProcessedBlock := o.startBlock
	blocksPerIteration := 10000

	for {
		o.extractor.OnStopCallback = func() {
			o.scanner.StopScanning()
		}

		go o.extractor.StartExtraction(
			lastProcessedBlock,
			o.threadsCount,
			lastProcessedBlock+blocksPerIteration,
			blocksPerIteration/o.threadsCount)
		o.scanner.StartScanning()

		fmt.Println("Orchestrator > Iteration finished. Scanned blocks are from", lastProcessedBlock, "to", o.extractor.LastProcessedBlock)
		lastProcessedBlock = o.extractor.LastProcessedBlock

		if len(o.scanner.Wallets) > 0 {
			o.storage.StoreSnapshot(lastProcessedBlock, o.scanner.Wallets)
		}
		o.storage.SaveTransactions(lastProcessedBlock, "processed-transactions", o.scanner.ProcessedTransactions)
		o.storage.SaveTransactions(lastProcessedBlock, "ignored-transactions", o.scanner.IgnoredTransactions)

		reinitializeOrchestrator(o, o.scanner.Wallets)
	}
}
