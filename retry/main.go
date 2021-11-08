package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type Effector func(context.Context) (string, error)

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (string, error) {
		for r := 0; ; r++ {
			res, err := effector(ctx)
			if err == nil || r >= retries {
				return res, err
			}

			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)

			<-time.After(delay)
		}
	}
}

func Test(context.Context) (string, error) {
	return "", errors.New("Failed")
}

func main() {
	ctx := context.Background()
	wrapped := Retry(Test, 3, 1*time.Second)
	res, err := wrapped(ctx)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(res)
	}
}
