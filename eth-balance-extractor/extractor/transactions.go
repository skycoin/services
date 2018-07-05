package extractor

import (
	"fmt"
	"strings"

	"github.com/golang-collections/go-datastructures/queue"
	"github.com/onrik/ethrpc"
)

type ExtractorStoppedCallback func()

// Extractor class
type Extractor struct {
	TransactionsQueue *queue.RingBuffer
	Stop              chan int
	IsStopped         bool
	IsDisposed        bool

	ContractAddress    string
	NodeAPI            string
	LastProcessedBlock int
	TransactionsLimit  int

	ExtractorStoppedCallback ExtractorStoppedCallback
}

// NewExtractor creates a new Extractor class
func NewExtractor(nodeAPI string, contractAddress string) *Extractor {
	return &Extractor{
		TransactionsQueue: queue.NewRingBuffer(100000),
		IsStopped:         false,
		IsDisposed:        false,

		ContractAddress:    strings.ToLower(contractAddress),
		NodeAPI:            nodeAPI,
		LastProcessedBlock: 0,
		TransactionsLimit:  0,
	}
}

func extraction(client *ethrpc.EthRPC, startBlock int, endBlock int, contractAddress string, queue *queue.RingBuffer, stopChannel chan int, id int) {
	fmt.Println("Started thread with id ", id)
	for i := startBlock; i <= endBlock; i++ {
		block, err := client.EthGetBlockByNumber(i, true)
		if block == nil {
			continue
		}
		if err != nil {
			fmt.Println(err)
		}
		for _, t := range block.Transactions {
			if strings.ToLower(t.To) == contractAddress {
				err := queue.Put(t)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	fmt.Println("Stopped thread ", id)
	stopChannel <- id
}

// StartExtraction start the extraction process
func (e *Extractor) StartExtraction(startBlock int, threadsCount int, blocksLimit int, blockPerThread int) {
	client := ethrpc.New(e.NodeAPI)

	threadID := 1

	e.Stop = make(chan int, threadsCount*2)

	for i := 0; i < threadsCount; i++ {
		start := startBlock + i*blockPerThread
		go extraction(client, start, start+blockPerThread, e.ContractAddress, e.TransactionsQueue, e.Stop, threadID)
		threadID++
	}
	lastProcessedBlock := startBlock + threadsCount*blockPerThread

	for {
		msg := <-e.Stop

		if msg == 0 {
			fmt.Println("Stop message has been received")
			e.LastProcessedBlock = lastProcessedBlock
			e.IsStopped = true
			e.TransactionsQueue.Put(nil)
		} else {
			threadsCount--
		}
		fmt.Println("Active threads: ", threadsCount)

		if threadsCount == 0 {
			e.IsDisposed = true
			e.LastProcessedBlock = lastProcessedBlock
			if e.ExtractorStoppedCallback != nil {
				e.ExtractorStoppedCallback()
			}
			return
		}

		if e.IsStopped {
			continue
		}

		fmt.Println("Transactions buffer length: ", e.TransactionsQueue.Len())

		if lastProcessedBlock >= blocksLimit {
			fmt.Println("Blocks limit has been reached ", lastProcessedBlock)
			go e.StopExtraction()
			continue
		}

		if threadsCount == 1 || e.TransactionsQueue.Len() < 50000 {
			fmt.Println("Started new thread. Starting block is ", lastProcessedBlock)
			go extraction(client, lastProcessedBlock, lastProcessedBlock+blockPerThread, e.ContractAddress, e.TransactionsQueue, e.Stop, threadID)
			threadsCount++
			threadID++
		}

		lastProcessedBlock = lastProcessedBlock + blockPerThread
	}
}

// StopExtraction stops the extraction process
func (e *Extractor) StopExtraction() {
	e.Stop <- 0

	for {
		if e.IsDisposed {
			fmt.Println("Extractor is stopped")
			if e.ExtractorStoppedCallback != nil {
				e.ExtractorStoppedCallback()
			}
			return
		}
	}
}
