package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func ContextWorkerPool() {
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
