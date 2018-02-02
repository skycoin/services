package main

import (
	"io"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/wallet"
)

var (
	LogWriter io.Writer

	STEP *Step
)

func init() {
	// for getting random Step.To addresses
	rand.Seed(time.Now().Unix())

	var err error
	if err = initFlags(); err != nil {
		log.Fatalln(err)
	}

	// create and check webrpc client
	c := &webrpc.Client{
		Addr: "localhost:6430",
	}
	if s, err := c.GetStatus(); err != nil {
		log.Fatalln(err)
	} else if !s.Running {
		log.Fatalln("node is not running on localhost:6430")
	}

	// create wallet
	w, err := wallet.NewWallet("load-testing", wallet.Options{
		Coin:  wallet.CoinTypeSkycoin,
		Label: "load-testing",
		Seed:  *SEED,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// initialize STEP for testing
	if STEP, err = NewStep(
		c,
		w,
		log.New(LogWriter, "[LOAD-TESTING] ", log.LstdFlags),
		w.GenerateAddresses(uint64(*N)),
	); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	// for graceful shutdown / cleanup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

transactions:
	for {
		result := STEP.Run()

		// log/record result

		println(result)

		select {
		case <-time.After(time.Second * time.Duration(*WAIT)):
			continue
		case <-stop:
			break transactions
		}
	}

	// finish transaction

	if *CLEANUP {
		// cleanup
	}
}
