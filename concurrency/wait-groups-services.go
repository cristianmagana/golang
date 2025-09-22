package main

import (
	"fmt"
	"sync"
	"time"
)

// Service function that sends status updates at specified intervals
func service(name string, interval time.Duration, updateChan chan<- string, wg *sync.WaitGroup, done <-chan bool) {
	defer wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C: // Wait for the ticker to "tick"
			// When ticker ticks, try to send a message
			select {
			case updateChan <- fmt.Sprintf("%s: OK", name):
			case <-done:
				return
			}
		case <-done:
			return
		}
	}
}

func WaitGroupServices() {
	fmt.Println("\nMulti-service monitoring system")

	var wg sync.WaitGroup

	wg.Add(3)

	// Create channels for each service
	dbChan := make(chan string, 10)
	apiChan := make(chan string, 10)
	cacheChan := make(chan string, 10)
	done := make(chan bool)

	// Start services with different intervals
	go service("DB", time.Second*2, dbChan, &wg, done)
	go service("API", time.Second*1, apiChan, &wg, done)
	go service("CACHE", time.Second*3, cacheChan, &wg, done)

	// Counters for each service
	dbCounter, apiCounter, cacheCounter := 0, 0, 0

	// Monitor for 10 seconds
	start := time.Now()
	timeout := time.Millisecond * 500

	for time.Since(start) < 10*time.Second {
		select {
		case msg := <-dbChan:
			fmt.Println(msg)
			dbCounter++
		case msg := <-apiChan:
			fmt.Println(msg)
			apiCounter++

		case msg := <-cacheChan:
			fmt.Println(msg)
			cacheCounter++
		case <-time.After(timeout):
			fmt.Println("Waiting for updates...")
		}
	}

	// Stop all services
	close(done)

	// Wait for all services to finish
	wg.Wait()

	// Report results
	fmt.Printf("\nMonitoring completed. Updates received:\n")
	fmt.Printf("Database service: %d updates\n", dbCounter)
	fmt.Printf("API service: %d updates\n", apiCounter)
	fmt.Printf("Cache service: %d updates\n", cacheCounter)
	fmt.Printf("Total updates: %d\n", dbCounter+apiCounter+cacheCounter)
}
