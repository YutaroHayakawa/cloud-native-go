package main

import (
	"fmt"
	"sync"
	"time"
)

func Funnel(sources ...chan int) chan int {
	dest := make(chan int)

	var wg sync.WaitGroup

	wg.Add(len(sources))

	for _, ch := range sources {
		go func(c <-chan int) {
			defer wg.Done()

			for n := range c {
				dest <- n
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(dest)
	}()

	return dest
}

func main() {
	var sources []chan int

	for i := 0; i < 3; i++ {
		ch := make(chan int)
		go func(id int, c chan int) {
			for {
				c <- id
				time.Sleep(1 * time.Second)
			}
		}(i, ch)
		sources = append(sources, ch)
	}

	dest := Funnel(sources...)

	for n := range dest {
		fmt.Println(n)
	}
}
