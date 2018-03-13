package btc

import (
	"log"
	"time"
)

/*
	Circuit breaker implementation https://martinfowler.com/bliki/CircuitBreaker.html
 	Intended to do actions with remote services that may go down for a while.
*/
type CircuitBreaker struct {
	success  func(string) (interface{}, error)
	fallback func(string) (interface{}, error)

	isOpen      uint32
	openTimeout time.Duration
	retryCount  uint
}

func NewCurcuitBreaker(success, fallback func(string) (interface{}, error), openTimeout time.Duration, retryCount uint) *CircuitBreaker {
	return &CircuitBreaker{
		success:  success,
		fallback: fallback,

		isOpen:      0,
		openTimeout: openTimeout,
		retryCount:  retryCount,
	}
}

func (c *CircuitBreaker) Do(arg string) (interface{}, error) {
	// If breaker is open - get info from block explorer
	if c.isOpen == 1 {
		result, err := c.fallback(arg)

		if err != nil {
			return 0, err
		}

		return result, nil
	}

	var i uint = 0

	result, err := c.success(arg)
	if err != nil {
		log.Printf("Get result from node returned error %s", err.Error())
	}

	for i < c.retryCount && err != nil {
		if err != nil {
			log.Printf("Get result from node returned error %s", err.Error())
		}

		result, err = c.success(arg)

		if err != nil {
			time.Sleep(time.Millisecond * time.Duration(1<<i) * 100)
		}
		i++
	}

	if err != nil {
		c.isOpen = 1

		go func() {
			time.Sleep(c.openTimeout)
			// This assignment is atomic since on 64-bit platforms this operation is atomic
			c.isOpen = 0
		}()

		result, err := c.fallback(arg)

		if err != nil {
			return 0.0, err
		}

		return result, nil
	}

	return result, nil
}

func (c *CircuitBreaker) IsOpen() bool {
	return c.isOpen == 1
}
