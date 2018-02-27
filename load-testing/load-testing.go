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
	// determined by flags, will contain either os.Stdout and a log file, or
	// just os.Stdout
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
		Addr: *NODE,
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

	addrs, err := w.GenerateAddresses(uint64(*N))
	if err != nil {
		log.Fatalln(err)
	}

	// initialize STEP for testing
	if STEP, err = NewStep(
		c,
		w,
		log.New(LogWriter, "", 0),
		addrs,
	); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	// for graceful shutdown / cleanup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	var (
		res *StepResult
		err error
	)

	summary := NewSummary(STEP.Addrs[0])

transactions:
	for {
		res, err = STEP.Run()
		if err != nil {
			panic(err)
		}

		summary.Add(res)

		if *LOG_TXS {
			STEP.Logger.Println(res.String())
		}

		select {
		case <-time.After(time.Second * time.Duration(*WAIT)):
			continue
		case <-stop:
			break transactions
		}
	}

	// stop the timer
	summary.Stop()

	if *CLEANUP {
		println("cleaning up...")
		res, err = STEP.Cleanup()
		STEP.Logger.Println(res.String())
	}

	if *LOG_SUM {
		STEP.Logger.Println(summary.String())
	}
}
