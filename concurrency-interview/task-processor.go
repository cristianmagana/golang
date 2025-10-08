package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Task represents a unit of work
type Task struct {
	ID      string
	Payload string
	Retries int
}

// Result represents the outcome of task processing
type Result struct {
	TaskID  string
	Success bool
	Error   error
}

// --------------------
// TaskProcessor
// --------------------

type TaskProcessor struct {
	WorkerCount int
	Tasks       chan Task
	Results     chan Result
	done        chan struct{}
	metrics     *Metrics
	wg          sync.WaitGroup
}

// NewTaskProcessor initializes the processor with given worker count
func NewTaskProcessor(workerCount int) *TaskProcessor {
	return &TaskProcessor{
		WorkerCount: workerCount,
		Tasks:       make(chan Task, 100),
		Results:     make(chan Result, 100),
		done:        make(chan struct{}),
		metrics:     &Metrics{},
	}
}

// Submit adds a task to the queue
func (tp *TaskProcessor) Submit(task Task) error {
	select {
	case <-tp.done:
		return fmt.Errorf("cannot submit, processor stopped")
	default:
		tp.metrics.IncrementTotal()
		tp.Tasks <- task
		return nil
	}
}

func (tp *TaskProcessor) StartWoRetry(ctx context.Context) {
	for i := 0; i < tp.WorkerCount; i++ {
		tp.wg.Add(1)
		go func(workerID int) {
			defer tp.wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case <-tp.done:
					return
				case task, ok := <-tp.Tasks:
					if !ok {
						return // channel closed
					}

					// Simulate per-task timeout
					taskCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
					err := processTask(taskCtx, task)
					cancel()

					success := err == nil
					tp.Results <- Result{
						TaskID:  task.ID,
						Success: success,
						Error:   err,
					}
				}
			}
		}(i)
	}
}

// Start begins the worker pool
func (tp *TaskProcessor) Start(ctx context.Context) {
	for i := 0; i < tp.WorkerCount; i++ {
		tp.wg.Add(1)
		go func(workerID int) {
			defer tp.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case <-tp.done:
					return
				case task, ok := <-tp.Tasks:
					if !ok {
						return // channel closed
					}

					// Per-task timeout context
					taskCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
					err := processTask(taskCtx, task)
					cancel()

					// Handle retry
					if err != nil && task.Retries < 3 {
						task.Retries++
						tp.metrics.IncrementRetried()
						backoff := time.Duration(task.Retries) * 200 * time.Millisecond
						time.Sleep(backoff)
						tp.Tasks <- task // requeue
						continue
					}

					if err != nil {
						tp.metrics.IncrementFailed()
						tp.Results <- Result{TaskID: task.ID, Success: false, Error: err}
					} else {
						tp.metrics.IncrementSuccess()
						tp.Results <- Result{TaskID: task.ID, Success: true}
					}
				}
			}
		}(i)
	}
}

// Stop gracefully shuts down the processor
func (tp *TaskProcessor) Stop() {
	close(tp.done)    // signal stop
	close(tp.Tasks)   // stop accepting tasks
	tp.wg.Wait()      // wait for workers to finish
	close(tp.Results) // signal results done
}

// --------------------
// Metrics
// --------------------

type Metrics struct {
	mu           sync.Mutex
	TotalTasks   int
	SuccessTasks int
	FailedTasks  int
	RetriedTasks int
}

func (m *Metrics) IncrementTotal() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalTasks++
}

func (m *Metrics) IncrementSuccess() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SuccessTasks++
}

func (m *Metrics) IncrementFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.FailedTasks++
}

func (m *Metrics) IncrementRetried() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RetriedTasks++
}

func (m *Metrics) GetStats() (int, int, int, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.TotalTasks, m.SuccessTasks, m.FailedTasks, m.RetriedTasks
}

// --------------------
// Task simulation
// --------------------

func processTask(ctx context.Context, task Task) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Duration(100+task.Retries*50) * time.Millisecond):
		if time.Now().UnixNano()%3 == 0 { // 30% failure rate
			return fmt.Errorf("task %s failed", task.ID)
		}
		return nil
	}
}

// --------------------
// Main
// --------------------

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	processor := NewTaskProcessor(5)
	go processor.Start(ctx)

	for i := 0; i < 20; i++ {
		task := Task{
			ID:      fmt.Sprintf("task-%d", i),
			Payload: fmt.Sprintf("data-%d", i),
		}
		if err := processor.Submit(task); err != nil {
			fmt.Printf("Submit failed: %v\n", err)
		}
	}

	// Collect results asynchronously
	go func() {
		for res := range processor.Results {
			if res.Success {
				fmt.Printf("[SUCCESS] %s processed\n", res.TaskID)
			} else {
				fmt.Printf("[FAIL] %s error: %v\n", res.TaskID, res.Error)
			}
		}
	}()

	time.Sleep(5 * time.Second)
	processor.Stop()

	total, success, failed, retried := processor.metrics.GetStats()
	fmt.Printf("\n--- METRICS ---\nTotal: %d | Success: %d | Failed: %d | Retried: %d\n",
		total, success, failed, retried)
	fmt.Println("All tasks processed")
}
