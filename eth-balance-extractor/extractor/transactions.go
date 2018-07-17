package extractor

import (
	"fmt"

	"github.com/golang-collections/go-datastructures/queue"
	"github.com/onrik/ethrpc"
)

// OnStopCallback represents a callback function that is invoked when extractor finished its work
type OnStopCallback func()

// Extractor class
type Extractor struct {
	TransactionsQueue *queue.RingBuffer
	Stop              chan int
	IsStopped         bool
	IsDisposed        bool

	NodeAPI            string
	LastProcessedBlock int
	TransactionsLimit  int

	OnStopCallback OnStopCallback
}

const RingBufferLength = 100000

// NewExtractor creates a new Extractor class
func NewExtractor(nodeAPI string) *Extractor {
	return &Extractor{
		TransactionsQueue: queue.NewRingBuffer(RingBufferLength),
		IsStopped:         false,
		IsDisposed:        false,

		NodeAPI:            nodeAPI,
		LastProcessedBlock: 0,
		TransactionsLimit:  0,
	}
}

func extraction(client *ethrpc.EthRPC, startBlock int, endBlock int, queue *queue.RingBuffer, stopChannel chan int, id int) {
	fmt.Println("Extractor > Started thread with id ", id)
	for i := startBlock; i <= endBlock; i++ {
		block, err := client.EthGetBlockByNumber(i, true)
		if err != nil {
			fmt.Println("Extractor > extraction", err)
			panic(err)
		}
		if block == nil {
			continue
		}
		if i%100 == 0 {
			fmt.Println("Extractor > thread", id, "processing block", i)
		}
		for _, t := range block.Transactions {
			err := queue.Put(t)
			if err != nil {
				fmt.Println("Extractor > extraction", err)
			}
		}
	}
	fmt.Println("Extractor > Stopped thread ", id)
	stopChannel <- id
}

// StartExtraction start the extraction process
func (e *Extractor) StartExtraction(startBlock int, threadsCount int, blocksLimit int, blockPerThread int) {
	client := ethrpc.New(e.NodeAPI)

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

		if msg == 0 {
			fmt.Println("Extractor > Stop message has been received")
			e.LastProcessedBlock = lastProcessedBlock
			e.IsStopped = true
			e.TransactionsQueue.Put(nil)
		} else {
			threadsCount--
		}
		fmt.Println("Extractor > Active threads: ", threadsCount)

		if threadsCount == 0 {
			e.IsDisposed = true
			e.LastProcessedBlock = lastProcessedBlock
			if e.OnStopCallback != nil {
				e.OnStopCallback()
			}
			return
		}

		if e.IsStopped {
			continue
		}

		fmt.Println("Extractor > Transactions queue length: ", e.TransactionsQueue.Len())

		if lastProcessedBlock >= blocksLimit && !e.IsStopped {
			fmt.Println("Extractor > Blocks limit has been reached ", lastProcessedBlock)
			e.IsStopped = true
			continue
		}

		if threadsCount == 1 || e.TransactionsQueue.Len() < (RingBufferLength/2) {
			fmt.Println("Extractor > Started new thread. Starting block is ", lastProcessedBlock)
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
			fmt.Println("Extractor > Extractor is stopped")
			if e.OnStopCallback != nil {
				e.OnStopCallback()
			}
			return
		}
	}
}
