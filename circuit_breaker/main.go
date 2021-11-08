package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Circuit func(context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock()

		d := consecutiveFailures - int(failureThreshold)

		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", fmt.Errorf("service unreachable, should try at %v", shouldRetryAt)
			}
		}

		m.RUnlock()

		response, err := circuit(ctx)

		m.Lock()
		defer m.Unlock()

		lastAttempt = time.Now()

		if err != nil {
			consecutiveFailures++
			return response, err
		}

		consecutiveFailures = 0

		return response, nil
	}
}

func FailureOperation(ctx context.Context) (string, error) {
	time.Sleep(1 * time.Second)
	return "", errors.New("Failed")
}

func main() {
	ctx := context.Background()
	breaker := Breaker(FailureOperation, 3)

	for i := 0; i < 4; i++ {
		_, err := breaker(ctx)
		if err != nil {
			fmt.Printf("[Failed] %v\n", err.Error())
		}
	}

	time.Sleep(5 * time.Second)

	fmt.Printf("Retrying at %v\n", time.Now())
	_, err := breaker(ctx)
	if err != nil {
		fmt.Printf("Got: %v\n", err.Error())
	}
}
