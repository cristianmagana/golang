package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func problem1_buggy_mutex() {
	counter := 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	start := time.Now()

	// Launch 10 goroutines, each incrementing counter 1000 times
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				counter++ // RACE CONDITION HERE
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Total elapsed time mutex: %d\n", time.Since(start).Microseconds())
	fmt.Printf("Problem 1 - Final counter: %d (expected 10000)\n", counter)
}

func problem1_buggy_atomic() {
	var counter int64
	var wg sync.WaitGroup

	start := time.Now()

	// Launch 10 goroutines, each incrementing counter 1000 times
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Total elapsed time atomic: %d\n", time.Since(start).Microseconds())
	fmt.Printf("Problem 1 - Final counter: %d (expected 10000)\n", counter)
}
