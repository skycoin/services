package main

import (
	"errors"
	"time"
)

type Duration struct {
	Start time.Time
	End   time.Time
}

type TransactionResult int

const (
	SUCCESS TransactionResult = iota
	FAIL
)

type Transaction struct {
	Duration *Duration
	Result   TransactionResult
}

type Summary struct {
	Duration     *Duration
	CleanupAddr  string
	Transactions map[string]*Transaction
}

func (t *Summary) Succeeded() int {
	count := 0

	for _, tx := range t.Transactions {
		if tx.Result == SUCCESS {
			count++
		}
	}

	return count
}

func (t *Summary) Failed() int {
	count := 0

	for _, tx := range t.Transactions {
		if tx.Result == FAIL {
			count++
		}
	}

	return count
}

func (t *Summary) Average() time.Duration {
	var (
		count int64 = int64(len(t.Transactions))
		total int64
	)

	for _, tx := range t.Transactions {
		total = total + int64(tx.Duration.End.Sub(tx.Duration.Start))
	}

	return time.Duration(total / count)
}

func (t *Summary) Minimum() time.Duration {
	var (
		minimum  int64 = -1
		duration int64
	)

	for _, tx := range t.Transactions {
		duration = int64(tx.Duration.End.Sub(tx.Duration.Start))

		if duration < minimum || minimum == -1 {
			minimum = duration
		}
	}

	return time.Duration(minimum)
}

func (t *Summary) Maximum() time.Duration {
	var (
		maximum  int64 = -1
		duration int64
	)

	for _, tx := range t.Transactions {
		duration = int64(tx.Duration.End.Sub(tx.Duration.Start))

		if duration > maximum || maximum == -1 {
			maximum = duration
		}
	}

	return time.Duration(maximum)
}

func (t *Summary) Start(txId string) error {
	if _, exists := t.Transactions[txId]; exists {
		return errors.New("timer for that txId has already been started")
	}

	t.Transactions[txId] = &Transaction{
		Duration: &Duration{
			Start: time.Now(),
		},
	}

	return nil
}

func (t *Summary) End(txId string) error {
	if _, exists := t.Transactions[txId]; !exists {
		return errors.New("timer for that txId doesn't exist")
	}

	t.Transactions[txId].Duration.End = time.Now()

	return nil
}
