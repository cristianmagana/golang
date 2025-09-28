package main

import (
	"fmt"
	"sync"
	"time"
)

func service(name string, interval time.Duration, updateChan chan<- string, wg *sync.WaitGroup, done <-chan bool) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case updateChan <- fmt.Sprintf("OK: %s\n", name):
			case <-done:
				return
			}
		case <-done:
			return
		}
	}
}

func WaitGroupServices() {

	// Initialize a waitgroup
	var wg sync.WaitGroup

	wg.Add(3)

	// Instantiate service and done channels
	dbChan := make(chan string, 10)
	apiChan := make(chan string, 10)
	cacheChan := make(chan string, 10)
	done := make(chan bool)

	dbCount, apiCount, cacheCount := 0, 0, 0

	go service("DB", 2*time.Second, dbChan, &wg, done)
	go service("API", 1*time.Second, apiChan, &wg, done)
	go service("CACHE", 3*time.Second, cacheChan, &wg, done)

	start := time.Now()
	timeout := time.Millisecond * 500

	for time.Since(start) < 10*time.Second {
		select {
		case msg := <-dbChan:
			fmt.Printf("%s\n", msg)
			dbCount++
		case msg := <-apiChan:
			fmt.Printf("%s\n", msg)
			apiCount++
		case msg := <-cacheChan:
			fmt.Printf("%s\n", msg)
			cacheCount++
		case <-time.After(timeout):
			fmt.Println("Waiting for updates...")
		}
	}

	close(done)

	wg.Wait()

	fmt.Printf("\nMonitoring completed. Updates received:\n")
	fmt.Printf("Database service: %d updates\n", dbCount)
	fmt.Printf("API service: %d updates\n", apiCount)
	fmt.Printf("Cache service: %d updates\n", cacheCount)
	fmt.Printf("Total updates: %d\n", dbCount+apiCount+cacheCount)
}
