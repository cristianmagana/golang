package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job represents a unit of work
type Job struct {
	ID      int
	Payload string
	Retries int
}

// JobResult represents the outcome of processing a job
type JobResult struct {
	JobID   int
	Success bool
	Error   error
}

// processJob simulates processing a job
// DO NOT MODIFY THIS FUNCTION
func processJob(job Job) error {
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
func startWorkerPool(ctx context.Context, jobs <-chan Job, results chan<- JobResult, numWorkers int, maxRetries int) {
	// Your implementation here
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					fmt.Printf("[Submitter] %d context exceeded, channel closed", id)
					return
				case job, ok := <-jobs:
					if !ok {
						fmt.Printf("error reading job %d from jobs channel, closing channel\n", id)
						return
					}
					fmt.Printf("[Worker %d] Received job %d\n", id, job.ID)

					for {
						err := processJob(job)
						if err == nil {
							results <- JobResult{JobID: job.ID, Success: true, Error: nil}
							break
						}

						job.Retries++
						fmt.Printf("[Worker %d] Job %d failed (retry %d/%d): %v\n",
							id, job.ID, job.Retries, maxRetries, err)

						if job.Retries > maxRetries {
							results <- JobResult{JobID: job.ID, Success: false, Error: err}
							break
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

// TODO: Implement this function
// This function should:
// - Send all jobs to the jobs channel
// - Close the jobs channel when done
func submitJobs(ctx context.Context, jobs chan<- Job, jobList []Job) {
	// Your implementation here

	defer close(jobs)

	for _, job := range jobList {
		select {
		case <-ctx.Done():
			fmt.Println("[Submitter] Context cancelled, stopping job submission")
			return
		case jobs <- job:
			fmt.Printf("[Submitter] Job %d submitted to channel\n", job.ID)
		}
	}
}

func main2() {
	// Create 50 jobs
	jobList := make([]Job, 50)
	for i := range jobList {
		jobList[i] = Job{
			ID:      i + 1,
			Payload: fmt.Sprintf("task-%d", i+1),
			Retries: 0,
		}
	}

	// Create channels
	jobs := make(chan Job, 10) // Buffered channel for jobs
	results := make(chan JobResult, 50)

	// Create context to handle timeouts and handling failed workers gracefully
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Start worker pool with 5 workers, max 3 retries per job
	numWorkers := 5
	maxRetries := 3

	start := time.Now()

	// TODO: Start workers and submit jobs
	// You need to coordinate these properly
	go startWorkerPool(ctx, jobs, results, numWorkers, maxRetries)

	go submitJobs(ctx, jobs, jobList)

	// TODO: Collect results
	// How do you know when all jobs are done?

	success, fail := 0, 0
	for res := range results {
		if res.Success {
			success++
		} else {
			fail++
		}
	}

	elapsed := time.Since(start)

	fmt.Printf("Processed all jobs in %v\n", elapsed)
	fmt.Printf("Processed %d jobs, with success: %d and errors: %d", len(results), success, fail)
}
