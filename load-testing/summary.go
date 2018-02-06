package main

import (
	"errors"
	"fmt"
	"time"
)

type Summary struct {
	Start        time.Time
	End          time.Time
	CleanupAddr  string
	Transactions map[string]*StepResult
}

func (s *Summary) String() string {
	var out string

	out = out + "SUMMARY:\n"

	out = out + fmt.Sprintf("- Start:        %v\n", s.Start)
	out = out + fmt.Sprintf("- End:          %v\n", s.End)
	out = out + fmt.Sprintf("- Duration:     %v\n", s.End.Sub(s.Start))
	out = out + fmt.Sprintf("- Cleanup:      %s\n", s.CleanupAddr)
	out = out + fmt.Sprintf("- Transactions: %d\n", s.All())
	out = out + fmt.Sprintf("     - Average: %v\n", s.Average())
	out = out + fmt.Sprintf("     - Minimum: %v\n", s.Minimum())
	out = out + fmt.Sprintf("     - Maximum: %v\n", s.Maximum())

	return out
}

func NewSummary(addr string) *Summary {
	return &Summary{
		Start:        time.Now(),
		CleanupAddr:  addr,
		Transactions: make(map[string]*StepResult, 0),
	}
}

func (s *Summary) Stop() {
	s.End = time.Now()
}

func (s *Summary) All() int {
	return len(s.Transactions)
}

func (s *Summary) Succeeded() int {
	count := 0

	for _, tx := range s.Transactions {
		if tx.Status.Confirmed {
			count++
		}
	}

	return count
}

func (s *Summary) Failed() int {
	count := 0

	for _, tx := range s.Transactions {
		if !tx.Status.Confirmed {
			count++
		}
	}

	return count
}

func (s *Summary) Average() time.Duration {
	var (
		count int64 = int64(len(s.Transactions))
		total int64
	)

	for _, tx := range s.Transactions {
		total = total + int64(tx.Duration)
	}

	return time.Duration(total / count)
}

func (s *Summary) Minimum() time.Duration {
	var (
		minimum  int64 = -1
		duration int64
	)

	for _, tx := range s.Transactions {
		duration = int64(tx.Duration)

		if duration < minimum || minimum == -1 {
			minimum = duration
		}
	}

	return time.Duration(minimum)
}

func (s *Summary) Maximum() time.Duration {
	var (
		maximum  int64 = -1
		duration int64
	)

	for _, tx := range s.Transactions {
		duration = int64(tx.Duration)

		if duration > maximum || maximum == -1 {
			maximum = duration
		}
	}

	return time.Duration(maximum)
}

func (s *Summary) Add(res *StepResult) error {
	if _, exists := s.Transactions[res.Id]; exists {
		return errors.New("timer for that txId has already been started")
	}

	s.Transactions[res.Id] = res

	return nil
}
