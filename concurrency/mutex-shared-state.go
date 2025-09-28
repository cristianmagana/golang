package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func MutexSharedState() {
	fmt.Printf("Starting Mutex Shared State Exercise\n")

	counter := &Counter{}
	var wg sync.WaitGroup

	// 100 goroutines incrementing 100 times each
	for range 100 {
		wg.Go(func() {
			for range 100 {
				counter.Increment()
			}
		})
	}
	wg.Wait()

	fmt.Printf("Final counter value: %d\n", counter.Value())
}
