package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// =============================================================================
// 1. CONTEXT & CANCELLATION - Why It's Critical
// =============================================================================

// PROBLEM: Without context, operations can run forever and waste resources
func badExampleWithoutContext() {
	fmt.Println("=== BAD: Without Context ===")

	// This could hang forever if the operation takes too long
	result := make(chan string, 1)

	go func() {
		// Simulate slow database query or API call
		time.Sleep(5 * time.Second)
		result <- "Operation completed"
	}()

	// Main thread blocks forever - no way to cancel!
	fmt.Println(<-result)
}

// SOLUTION: With context, we can set timeouts and cancel operations
func goodExampleWithContext() {
	fmt.Println("\n=== GOOD: With Context ===")

	// Set a 2-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result := make(chan string, 1)

	go func() {
		// Simulate slow operation
		time.Sleep(5 * time.Second)
		result <- "Operation completed"
	}()

	select {
	case res := <-result:
		fmt.Println("✅", res)
	case <-ctx.Done():
		fmt.Println("❌ Operation timed out:", ctx.Err())
	}
}

// REAL-WORLD USE CASES for Context:
// 1. HTTP request timeouts
// 2. Database query timeouts
// 3. User cancels long-running operation
// 4. Service shutdown - cancel all ongoing work
// 5. Request tracing in microservices

func httpClientWithTimeout() {
	fmt.Println("\n=== HTTP Client with Timeout ===")

	// Create HTTP client with 1-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", "http://httpbin.org/delay/3", nil)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("❌ Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("✅ Response: %s\n", resp.Status)
}

// =============================================================================
// 2. GRACEFUL SHUTDOWN - Why It's Essential in Production
// =============================================================================

// PROBLEM: Abrupt shutdown can cause:
// - Data corruption
// - Lost requests
// - Incomplete transactions
// - Resource leaks

// SOLUTION: Graceful shutdown gives time to:
// - Complete current requests
// - Save state to disk
// - Close database connections
// - Send final metrics

func gracefulShutdownWebServer() {
	fmt.Println("\n=== Graceful Shutdown Web Server ===")

	// Create HTTP server
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate some work
			time.Sleep(2 * time.Second)
			fmt.Fprintf(w, "Request processed at %v", time.Now())
		}),
	}

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		fmt.Println("🚀 Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ Server error: %v\n", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	fmt.Println("📢 Shutdown signal received")

	// Give server 5 seconds to finish current requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("❌ Server forced shutdown: %v\n", err)
	} else {
		fmt.Println("✅ Server gracefully stopped")
	}
}

// WHEN TO USE GRACEFUL SHUTDOWN:
// - Web servers
// - Database connections
// - Message queue consumers
// - File processing services
// - Any long-running service

func gracefulShutdownDataProcessor() {
	fmt.Println("\n=== Graceful Shutdown Data Processor ===")

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Simulate data processing service
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer fmt.Println("✅ Data processor stopped")

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("📊 Processing batch...")
				// Simulate processing work
				time.Sleep(200 * time.Millisecond)

			case <-ctx.Done():
				fmt.Println("🛑 Finishing current batch before shutdown...")
				// Complete current work
				time.Sleep(300 * time.Millisecond)
				fmt.Println("💾 Saved state to disk")
				return
			}
		}
	}()

	// Simulate shutdown after 3 seconds
	time.Sleep(3 * time.Second)
	fmt.Println("📢 Initiating graceful shutdown...")
	cancel()

	wg.Wait()
	fmt.Println("🏁 All services stopped cleanly")
}

// =============================================================================
// 3. WORKER POOL WITH CONTEXT - Why It's Powerful
// =============================================================================

// PROBLEMS WITHOUT CONTEXT IN WORKER POOLS:
// - Workers can't be stopped
// - Memory leaks from abandoned goroutines
// - No way to handle timeouts
// - Difficult to coordinate shutdown

// SOLUTIONS WITH CONTEXT:
// - Coordinated cancellation
// - Timeout handling
// - Resource cleanup
// - Graceful worker shutdown

func workerPoolRealisticExample() {
	fmt.Println("\n=== Realistic Worker Pool with Context ===")

	// Context with 10-second overall timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Job types
	type Job struct {
		ID       int
		Data     string
		Priority int
	}

	type Result struct {
		JobID  int
		Output string
		Error  error
		Worker int
	}

	jobs := make(chan Job, 100)
	results := make(chan Result, 100)

	// Start worker pool
	var wg sync.WaitGroup
	numWorkers := 3

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			defer fmt.Printf("🔧 Worker %d shutdown\n", workerID)

			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						return // Jobs channel closed
					}

					// Process job with timeout check
					processStart := time.Now()

					// Simulate work that respects context
					workDone := make(chan Result, 1)
					go func() {
						// Simulate processing time
						processingTime := time.Duration(job.Priority) * time.Millisecond * 500
						time.Sleep(processingTime)

						workDone <- Result{
							JobID:  job.ID,
							Output: fmt.Sprintf("Processed '%s'", job.Data),
							Worker: workerID,
						}
					}()

					// Wait for work or context cancellation
					select {
					case result := <-workDone:
						result.Worker = workerID
						results <- result
						fmt.Printf("⚡ Worker %d completed job %d in %v\n",
							workerID, job.ID, time.Since(processStart))

					case <-ctx.Done():
						results <- Result{
							JobID:  job.ID,
							Error:  fmt.Errorf("job %d cancelled: %v", job.ID, ctx.Err()),
							Worker: workerID,
						}
						return // Exit worker
					}

				case <-ctx.Done():
					fmt.Printf("🛑 Worker %d received cancellation\n", workerID)
					return
				}
			}
		}(w)
	}

	// Job producer
	go func() {
		defer close(jobs)

		jobData := []struct {
			data     string
			priority int
		}{
			{"urgent-task", 1}, // Fast job
			{"normal-task", 2}, // Medium job
			{"slow-task", 4},   // Slow job
			{"batch-job", 3},   // Medium job
			{"cleanup-job", 1}, // Fast job
		}

		for i, job := range jobData {
			select {
			case jobs <- Job{ID: i + 1, Data: job.data, Priority: job.priority}:
				fmt.Printf("📨 Sent job %d: %s\n", i+1, job.data)
			case <-ctx.Done():
				fmt.Println("📨 Job producer cancelled")
				return
			}
		}
	}()

	// Result collector
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and display results
	for result := range results {
		if result.Error != nil {
			fmt.Printf("❌ %v\n", result.Error)
		} else {
			fmt.Printf("✅ Job %d: %s (by worker %d)\n",
				result.JobID, result.Output, result.Worker)
		}
	}

	fmt.Println("🏁 All work completed or cancelled")
}

// =============================================================================
// WHEN TO USE EACH PATTERN
// =============================================================================

func whenToUseWhichPattern() {
	fmt.Println("\n=== WHEN TO USE WHICH PATTERN ===")

	fmt.Println(`
🎯 CONTEXT & CANCELLATION:
   ✅ Use when: API calls, database queries, file I/O, any blocking operation
   ✅ Types:
      - WithTimeout: "This operation should not take more than X seconds"
      - WithDeadline: "This must complete by 3:00 PM"
      - WithCancel: "User can cancel this anytime"
      - WithValue: "Pass request ID through call chain"

🛑 GRACEFUL SHUTDOWN:
   ✅ Use when: Web servers, background services, data processors
   ✅ Benefits: No data loss, complete current requests, clean resource cleanup
   ✅ Pattern: Listen for signals (SIGINT/SIGTERM), cancel context, wait for completion

⚡ WORKER POOL WITH CONTEXT:
   ✅ Use when: Processing many similar tasks, need parallelism + control
   ✅ Benefits: Controlled concurrency, timeout handling, coordinated shutdown
   ✅ Examples: Image processing, data transformation, API fanout calls
	`)
}

func ContextCases() {
	// Uncomment to run specific examples
	goodExampleWithContext()
	time.Sleep(1 * time.Second)

	// httpClientWithTimeout() // Requires network

	gracefulShutdownDataProcessor()
	time.Sleep(1 * time.Second)

	workerPoolRealisticExample()

	whenToUseWhichPattern()

	// gracefulShutdownWebServer() // Uncomment to test web server
}
