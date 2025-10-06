package main

import (
	"fmt"
	"sync"
	"time"
)

func problem2_buggy_map_rwmutext() {
	cache := make(map[string]int)
	var wg sync.WaitGroup
	var mu sync.RWMutex

	start := time.Now()

	for i := range 5 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range 100 {
				mu.Lock()
				key := fmt.Sprintf("key_%d", j)
				cache[key] = id
				mu.Unlock()
			}
		}(i)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			mu.RLock()
			_ = cache["key_50"]
			mu.RUnlock()
		}
	}()

	wg.Wait()
	fmt.Printf("Total elapsed time map rwmutex: %d\n", time.Since(start).Microseconds())
	fmt.Printf("Problem 2 - Map entries: %d\n", len(cache))
}

func problem2_buggy_map_syncmap() {
	var cache sync.Map
	var wg sync.WaitGroup

	start := time.Now()

	// Multiple goroutines writing
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key_%d", j)
				cache.Store(key, id)
			}
		}(i)
	}

	// Goroutine reading
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			cache.Load("key_50")
		}
	}()

	wg.Wait()
	count := 0
	cache.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	fmt.Printf("Total elapsed time syncmap: %d\n", time.Since(start).Microseconds())
	fmt.Printf("Problem 2 - Map entries: %d\n", count)
}
