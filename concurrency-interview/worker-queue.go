package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job represents a unit of work
type Job_ struct {
	ID      int
	Payload string
	Retries int
}

// JobResult represents the outcome of processing a job
type JobResult_ struct {
	JobID   int
	Success bool
	Error   error
}

// processJob simulates processing a job
// DO NOT MODIFY THIS FUNCTION
func processJob_(job Job) error {
	// Simulate variable processing time (50-300ms)
	processingTime := 50 + rand.Intn(250)
	time.Sleep(time.Duration(processingTime) * time.Millisecond)

	// Simulate 20% failure rate
	if rand.Intn(100) < 20 {
		return fmt.Errorf("processing failed for job %d", job.ID)
	}

	return nil
}

// TODO: Implement this function
// Create a worker pool that:
// - Spawns numWorkers goroutines
// - Each worker consumes jobs from the jobs channel
// - Failed jobs should be retried up to maxRetries times
// - Results (success or final failure) should be sent to results channel
// - Workers should stop gracefully when jobs channel is closed
func startWorkerPool_(ctx context.Context, jobs <-chan Job, results chan<- JobResult, numWorkers int, maxRetries int) {
	// Your implementation here
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("[Worker Pool] Worker %d context canceled, shutting down\n", id)
					return
				case job, ok := <-jobs:
					if !ok {
						fmt.Printf("[Worker %d] Worker channel shutting down\n", id)
						return
					}
					// Process with retries
					var err error
					for attempt := 0; attempt <= maxRetries; attempt++ {
						err = processJobsWithTimeout_(ctx, job)
						if err == nil {
							results <- JobResult{JobID: job.ID, Success: true, Error: nil}
							break
						}
						fmt.Printf("[Worker %d] Job %d failed (attempt %d/%d): %v\n", id, job.ID, attempt+1, maxRetries, err)
						if attempt == maxRetries {
							results <- JobResult{JobID: job.ID, Success: false, Error: err}
						}
					}
				}
			}
		}(i + 1)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
}

func processJobsWithTimeout_(ctx context.Context, job Job) error {
	done := make(chan error, 1)

	go func() {
		done <- processJob(job)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("job %d timed out", job.ID)
	}
}

// TODO: Implement this function
// This function should:
// - Send all jobs to the jobs channel
// - Close the jobs channel when done
func submitJobs_(ctx context.Context, jobs chan<- Job, jobList []Job) {
	defer close(jobs)
	for _, job := range jobList {
		select {
		case jobs <- job:
			fmt.Printf("[Submitter] Submitted job %d onto queue\n", job.ID)
		case <-ctx.Done():
			fmt.Println("[Submitter] Context cancelled")
			return
		}
	}
	fmt.Printf("[Submitter] All %d jobs have been submitted", len(jobList))
}

func WorkerQueue() {
	// Create 50 jobs
	jobList := make([]Job, 50)
	for i := range jobList {
		jobList[i] = Job{
			ID:      i + 1,
			Payload: fmt.Sprintf("task-%d", i+1),
			Retries: 0,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Create channels
	jobs := make(chan Job, 10) // Buffered channel for jobs
	results := make(chan JobResult, 50)

	// Start worker pool with 5 workers, max 3 retries per job
	numWorkers := 5
	maxRetries := 3

	start := time.Now()

	// TODO: Start workers and submit jobs
	// You need to coordinate these properly

	go startWorkerPool_(ctx, jobs, results, numWorkers, maxRetries)

	go submitJobs_(ctx, jobs, jobList)

	// TODO: Collect results
	// How do you know when all jobs are done?

	success, fail := 0, 0
	for result := range results {
		if result.Success {
			success++
		} else {
			fail++
		}
	}

	elapsed := time.Since(start)

	fmt.Printf("Processed all jobs in %v\n", elapsed)
	fmt.Printf("✅ Done. Success: %d, Failures: %d\n", success, fail)
}
