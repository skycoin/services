package extractor

import (
	"fmt"
)

// Orchestrator represents an orchestrator of transactions extractor and wallets scanner
type Orchestrator struct {
	extractor *Extractor
	scanner   *WalletScanner
	storage   *Storage

	startBlock   int
	threadsCount int
}

// NewOrchestrator creates a new instance of the Orchestrator
func NewOrchestrator(startBlock int, threadsCount int) *Orchestrator {
	e := NewExtractor()
	return &Orchestrator{
		extractor: e,
		scanner:   NewWalletScanner(e.TransactionsQueue),
		storage:   NewStorage("./tmp"),

		startBlock:   startBlock,
		threadsCount: threadsCount,
	}
}

func reinitializeOrchestrator(o *Orchestrator, wallets map[string]*Wallet) {
	e := NewExtractor()
	o.extractor = e
	o.scanner = NewWalletScanner(e.TransactionsQueue)
	o.scanner.Wallets = wallets
}

// StartScanning starts scanning process
func (o *Orchestrator) StartScanning() {
	lastProcessedBlock := o.startBlock
	blocksPerIteration := 100

	for {
		o.extractor.ExtractorStoppedCallback = func() {
			o.scanner.StopScanning()
		}

		go o.extractor.StartExtraction(
			lastProcessedBlock,
			o.threadsCount,
			lastProcessedBlock+blocksPerIteration,
			blocksPerIteration/o.threadsCount)
		o.scanner.StartScanning()

		fmt.Println("Iteration finished. Scanned blocks are from", lastProcessedBlock, "to", o.extractor.LastProcessedBlock)
		lastProcessedBlock = o.extractor.LastProcessedBlock

		o.storage.StoreSnapshot(lastProcessedBlock, o.scanner.Wallets)

		reinitializeOrchestrator(o, o.scanner.Wallets)
	}
}
