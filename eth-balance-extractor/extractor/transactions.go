package extractor

import (
	"fmt"

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

	LastProcessedBlock int
	TransactionsLimit  int

	ExtractorStoppedCallback ExtractorStoppedCallback
}

// NewExtractor creates a new Extractor class
func NewExtractor() *Extractor {
	return &Extractor{
		TransactionsQueue:  queue.NewRingBuffer(100000),
		LastProcessedBlock: 0,
		IsStopped:          false,
		IsDisposed:         false,

		TransactionsLimit: 0,
	}
}

func extraction(client *ethrpc.EthRPC, startBlock int, endBlock int, queue *queue.RingBuffer, stopChannel chan int, id int) {
	fmt.Println("Started thread with id ", id)
	for i := startBlock; i <= endBlock; i++ {
		block, err := client.EthGetBlockByNumber(i, true)
		if err != nil {
			fmt.Println(err)
		}
		for _, t := range block.Transactions {
			err := queue.Put(t)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	fmt.Println("Stopped thread ", id)
	stopChannel <- id
}

// StartExtraction start the extraction process
func (e *Extractor) StartExtraction(startBlock int, threadsCount int, blocksLimit int, blockPerThread int) {
	client := ethrpc.New("http://127.0.0.1:8545")

	threadID := 1

	e.Stop = make(chan int, threadsCount*2)

	for i := 0; i < threadsCount; i++ {
		start := startBlock + i*blockPerThread
		go extraction(client, start, start+blockPerThread, e.TransactionsQueue, e.Stop, threadID)
		threadID++
	}
	lastProcessedBlock := startBlock + threadsCount*blockPerThread

	for {
		msg := <-e.Stop
		fmt.Println("Active threads: ", threadsCount)

		if msg == 0 {
			fmt.Println("Stop message has been received")
			e.LastProcessedBlock = lastProcessedBlock
			e.IsStopped = true
			e.TransactionsQueue.Put(nil)
		} else {
			threadsCount--
		}

		if threadsCount == 0 {
			e.IsDisposed = true
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
			go extraction(client, lastProcessedBlock, lastProcessedBlock+blockPerThread, e.TransactionsQueue, e.Stop, threadID)
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
