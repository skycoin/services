package btc

import (
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestCircuitBreakerDo(t *testing.T) {
	var success = "success"

	testData := []struct {
		success  func(string) (interface{}, error)
		fallback func(string) (interface{}, error)
		expected bool
	}{
		{
			success: func(arg string) (interface{}, error) {
				return success, nil
			},
			fallback: func(arg string) (interface{}, error) {
				return struct{}{}, nil
			},
			expected: true,
		},
		{
			success: func(arg string) (interface{}, error) {
				return struct{}{}, errors.New("error")
			},
			fallback: func(arg string) (interface{}, error) {
				return success, nil
			},
			expected: false,
		},
	}

	for _, test := range testData {
		circuitBreaker := CircuitBreaker{
			success:  test.success,
			fallback: test.fallback,

			isOpen:      0,
			openTimeout: time.Second * 10,
			retryCount:  3,
		}

		result, err := circuitBreaker.Do("")

		if test.expected {
			if err != nil {
				t.Errorf("Expected success got error %s", err.Error())
				return
			}

			r, _ := result.(string)

			if r != success {
				t.Errorf("Wrong result expected %s actual %s", success, r)
			}

			if circuitBreaker.IsOpen() {
				t.Errorf("Expected circuit breaker not to be open, actual %t", circuitBreaker.IsOpen())
			}
		} else {
			if !circuitBreaker.IsOpen() {
				t.Errorf("Expected circuit breaker to be open, actual %t", circuitBreaker.IsOpen())
			}
		}
	}
}
