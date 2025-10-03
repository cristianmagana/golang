package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 1. CONTEXT & CANCELLATION (Critical for distributed systems)
func contextExample() {
	fmt.Println("=== Context & Cancellation ===")

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			select {
			case <-time.After(3 * time.Second): // This would timeout
				fmt.Printf("Worker %d completed\n", id)
			case <-ctx.Done():
				fmt.Printf("Worker %d cancelled: %v\n", id, ctx.Err())
			}
		}(i)
	}

	wg.Wait()
}

// 2. ERROR HANDLING PATTERNS
type Result struct {
	Data string
	Err  error
}

func errorHandlingPattern() {
	fmt.Println("\n=== Error Handling Patterns ===")

	results := make(chan Result, 5)

	go func() {
		defer close(results)

		// Simulate operations that can fail
		operations := []struct {
			data       string
			shouldFail bool
		}{
			{"task1", false},
			{"task2", true},
			{"task3", false},
		}

		for _, op := range operations {
			if op.shouldFail {
				results <- Result{Err: fmt.Errorf("failed processing %s", op.data)}
			} else {
				results <- Result{Data: fmt.Sprintf("processed %s", op.data)}
			}
		}
	}()

	// Handle results with proper error checking
	for result := range results {
		if result.Err != nil {
			fmt.Printf("❌ Error: %v\n", result.Err)
			continue
		}
		fmt.Printf("✅ Success: %s\n", result.Data)
	}
}

// 3. GRACEFUL SHUTDOWN PATTERN
func gracefulShutdown() {
	fmt.Println("\n=== Graceful Shutdown Pattern ===")

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Simulate long-running service
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("🔄 Service working...")
			case <-ctx.Done():
				fmt.Println("🛑 Service shutting down gracefully...")
				// Cleanup work here
				time.Sleep(100 * time.Millisecond)
				fmt.Println("✅ Service shutdown complete")
				return
			}
		}
	}()

	// Simulate shutdown signal after 2 seconds
	time.Sleep(2 * time.Second)
	fmt.Println("📢 Shutdown signal received")
	cancel()

	wg.Wait()
	fmt.Println("🏁 All services stopped")
}

// 4. RATE LIMITING / THROTTLING
func rateLimitingExample() {
	fmt.Println("\n=== Rate Limiting ===")

	// Rate limiter: 2 operations per second
	rate := time.Second / 2
	limiter := time.NewTicker(rate)
	defer limiter.Stop()

	requests := []string{"req1", "req2", "req3", "req4", "req5"}

	for _, req := range requests {
		<-limiter.C // Wait for rate limit
		fmt.Printf("⚡ Processing %s at %s\n", req, time.Now().Format("15:04:05"))
	}
}

// 5. WORKER POOL WITH CONTEXT
func workerPoolWithContext() {
	fmt.Println("\n=== Worker Pool with Context ===")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	jobs := make(chan string, 10)
	results := make(chan string, 10)

	// Start workers
	var wg sync.WaitGroup
	numWorkers := 3

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						fmt.Printf("Worker %d: jobs channel closed\n", id)
						return
					}

					// Simulate work
					select {
					case <-time.After(500 * time.Millisecond):
						results <- fmt.Sprintf("Worker %d processed %s", id, job)
					case <-ctx.Done():
						fmt.Printf("Worker %d: context cancelled\n", id)
						return
					}

				case <-ctx.Done():
					fmt.Printf("Worker %d: context cancelled\n", id)
					return
				}
			}
		}(w)
	}

	// Send jobs
	go func() {
		defer close(jobs)
		for i := 1; i <= 5; i++ {
			select {
			case jobs <- fmt.Sprintf("job%d", i):
				fmt.Printf("📨 Sent job%d\n", i)
			case <-ctx.Done():
				fmt.Println("📨 Job sender cancelled")
				return
			}
		}
	}()

	// Close results when workers done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for result := range results {
		fmt.Printf("📬 %s\n", result)
	}
}

func ContextExamples() {
	contextExample()
	time.Sleep(500 * time.Millisecond)

	errorHandlingPattern()
	time.Sleep(500 * time.Millisecond)

	gracefulShutdown()
	time.Sleep(500 * time.Millisecond)

	rateLimitingExample()
	time.Sleep(500 * time.Millisecond)

	workerPoolWithContext()
}
