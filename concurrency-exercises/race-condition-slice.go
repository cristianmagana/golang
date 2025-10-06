package main

import (
	"fmt"
	"sync"
	"time"
)

func problem3_buggy_slice_mutext() {
	var results []int
	var wg sync.WaitGroup
	var mu sync.Mutex

	start := time.Now()

	for i := range 100000 {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			mu.Lock()
			results = append(results, val)
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	fmt.Printf("Total elapsed time slice mutex: %d\n", time.Since(start).Microseconds())
	fmt.Printf("Solution 3 (Mutex) - Results length: %d\n", len(results))

}
func problem3_buggy_slice_channels() {
	var results []int
	var wg sync.WaitGroup
	work := make(chan int, 1000)

	start := time.Now()

	for i := range 100000 {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			work <- val
		}(i)
	}

	go func() {
		wg.Wait()
		close(work)
	}()

	for i := range work {
		results = append(results, i)
	}

	wg.Wait()
	fmt.Printf("Total elapsed time slice channel: %d\n", time.Since(start).Microseconds())
	fmt.Printf("Solution 3 (Channels) - Results length: %d\n", len(results))

}
