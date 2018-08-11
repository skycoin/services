package extractor

import (
	"fmt"
	"strconv"
)

// Orchestrator represents an orchestrator of transactions extractor and wallets scanner
type Orchestrator struct {
	extractor *Extractor
	scanner   *WalletScanner
	storage   *Storage

	nodeAPI string

	startBlock   int
	threadsCount int
}

// NewOrchestrator creates a new instance of the Orchestrator
func NewOrchestrator(nodeAPI string, destDir string, startBlock int, threadsCount int) *Orchestrator {
	e := NewExtractor(nodeAPI)
	return &Orchestrator{
		extractor: e,
		scanner:   NewWalletScanner(nodeAPI, e.TransactionsQueue, nil),
		storage:   NewStorage(destDir),
		nodeAPI:   nodeAPI,

		startBlock:   startBlock,
		threadsCount: threadsCount,
	}
}

func reinitializeOrchestrator(o *Orchestrator, wallets map[string]*Wallet) {
	e := NewExtractor(o.nodeAPI)
	o.extractor = e
	o.scanner = NewWalletScanner(o.nodeAPI, e.TransactionsQueue, nil)
	o.scanner.Wallets = wallets
}

// StartScanning starts scanning process
func (o *Orchestrator) StartScanning(walletsFile string) {
	lastProcessedBlock := o.startBlock
	blocksPerIteration := 10000

	wallets := o.storage.LoadTransactionWallets(walletsFile)
	o.scanner.TransactionWallets = wallets
	o.scanner.RestoreKeys()

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

		o.storage.StoreSnapshot(strconv.Itoa(lastProcessedBlock), o.scanner.Wallets)

		reinitializeOrchestrator(o, o.scanner.Wallets)
	}
}
