package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

// Metrics to track application performance
type Metrics struct {
	RequestCount      int64 `json:"request_count"`
	ErrorCount        int64 `json:"error_count"`
	ActiveWorkers     int64 `json:"active_workers"`
	ProcessedJobs     int64 `json:"processed_jobs"`
	PendingJobs       int64 `json:"pending_jobs"`
	Goroutines        int   `json:"goroutines"`
	MemAllocMB        int   `json:"mem_alloc_mb"`
	MemSysMB          int   `json:"mem_sys_mb"`
	GCRuns            int   `json:"gc_runs"`
	LeakedGoroutines  int64 `json:"leaked_goroutines"`
	BlockedGoroutines int64 `json:"blocked_goroutines"`
}

// Job represents work to be done
type Job struct {
	ID        int           `json:"id"`
	Payload   string        `json:"payload"`
	Delay     time.Duration `json:"delay"`
	LeakSize  int           `json:"leak_size"`
	CreatedAt time.Time     `json:"created_at"`
}

// Result represents job completion
type Result struct {
	JobID       int           `json:"job_id"`
	Success     bool          `json:"success"`
	Duration    time.Duration `json:"duration"`
	CompletedAt time.Time     `json:"completed_at"`
}

// Application holds our application state
type Application struct {
	metrics     *Metrics
	jobQueue    chan Job
	resultQueue chan Result
	workerPool  *WorkerPool
	memoryLeaks [][]byte // Intentional memory leak for testing
	mu          sync.RWMutex
	shutdown    chan struct{}

	// Problem simulation controls
	enableMemoryLeak    bool
	enableGoroutineLeak bool
	enableDeadlock      bool
	enableMutexLeak     bool

	// Problematic components
	deadlockMutex      sync.Mutex
	leakyChannels      []chan int
	mutexLeakResources []*MutexLeakResource
}

// MutexLeakResource simulates a resource that acquires mutex but never releases
type MutexLeakResource struct {
	mu   sync.Mutex
	data string
	id   int
}

func (mlr *MutexLeakResource) LeakyOperation() {
	mlr.mu.Lock()
	// Intentionally never unlock - simulates mutex leak
	mlr.data = fmt.Sprintf("processing-%d", mlr.id)
	time.Sleep(100 * time.Millisecond)
	// Missing: mlr.mu.Unlock()
}

func (mlr *MutexLeakResource) ProblematicAccess() {
	// This will block forever if LeakyOperation was called
	mlr.mu.Lock()
	defer mlr.mu.Unlock()
	log.Printf("Accessing resource %d: %s", mlr.id, mlr.data)
}

// WorkerPool manages our worker goroutines
type WorkerPool struct {
	workers       int
	jobs          chan Job
	results       chan Result
	activeJobs    int64
	processedJobs int64
	quit          chan struct{}
	wg            sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int, jobs chan Job, results chan Result) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		jobs:    jobs,
		results: results,
		quit:    make(chan struct{}),
	}
}

// Start begins the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.quit)
	wp.wg.Wait()
}

// worker processes jobs from the job queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	log.Printf("Worker %d started", id)

	for {
		select {
		case job := <-wp.jobs:
			atomic.AddInt64(&wp.activeJobs, 1)
			start := time.Now()

			result := wp.processJob(job, id)

			atomic.AddInt64(&wp.activeJobs, -1)
			atomic.AddInt64(&wp.processedJobs, 1)

			result.Duration = time.Since(start)
			result.CompletedAt = time.Now()

			select {
			case wp.results <- result:
			case <-wp.quit:
				return
			}

		case <-wp.quit:
			log.Printf("Worker %d stopping", id)
			return
		}
	}
}

// processJob simulates work and potential issues
func (wp *WorkerPool) processJob(job Job, workerID int) Result {
	result := Result{
		JobID:   job.ID,
		Success: true,
	}

	// Simulate work delay
	if job.Delay > 0 {
		time.Sleep(job.Delay)
	}

	// Simulate CPU-intensive work
	sum := 0
	for i := 0; i < 1000000; i++ {
		sum += i
	}

	// Simulate memory allocation (potential leak)
	if job.LeakSize > 0 {
		_ = make([]byte, job.LeakSize)
	}

	// Random chance of failure to test error handling
	if rand.Float32() < 0.05 { // 5% failure rate
		result.Success = false
	}

	if job.ID%100 == 0 {
		log.Printf("Worker %d processed job %d (payload: %s)", workerID, job.ID, job.Payload)
	}

	return result
}

// NewApplication creates a new application instance
func NewApplication() *Application {
	jobQueue := make(chan Job, 1000)
	resultQueue := make(chan Result, 1000)

	app := &Application{
		metrics:     &Metrics{},
		jobQueue:    jobQueue,
		resultQueue: resultQueue,
		workerPool:  NewWorkerPool(10, jobQueue, resultQueue),
		shutdown:    make(chan struct{}),

		// Enable problems based on environment variables
		enableMemoryLeak:    os.Getenv("ENABLE_MEMORY_LEAK") == "true",
		enableGoroutineLeak: os.Getenv("ENABLE_GOROUTINE_LEAK") == "true",
		enableDeadlock:      os.Getenv("ENABLE_DEADLOCK") == "true",
		enableMutexLeak:     os.Getenv("ENABLE_MUTEX_LEAK") == "true",

		leakyChannels:      make([]chan int, 0),
		mutexLeakResources: make([]*MutexLeakResource, 0),
	}

	return app
}

// Start begins the application
func (app *Application) Start() {
	log.Println("Starting application...")
	log.Printf("Problem simulation - Memory Leak: %v, Goroutine Leak: %v, Deadlock: %v, Mutex Leak: %v",
		app.enableMemoryLeak, app.enableGoroutineLeak, app.enableDeadlock, app.enableMutexLeak)

	// Start worker pool
	app.workerPool.Start()

	// Start job generator
	go app.generateJobs()

	// Start result processor
	go app.processResults()

	// Start metrics updater
	go app.updateMetrics()

	// Start problem simulators
	if app.enableMemoryLeak {
		go app.simulateMemoryLeak()
	}

	if app.enableGoroutineLeak {
		go app.simulateGoroutineLeak()
	}

	if app.enableDeadlock {
		go app.simulateDeadlock()
	}

	if app.enableMutexLeak {
		go app.simulateMutexLeak()
	}

	log.Println("Application started")
}

// Stop gracefully shuts down the application
func (app *Application) Stop() {
	log.Println("Stopping application...")
	close(app.shutdown)
	app.workerPool.Stop()
	log.Println("Application stopped")
}

// generateJobs creates work for the system
func (app *Application) generateJobs() {
	jobID := 0
	ticker := time.NewTicker(100 * time.Millisecond) // 10 jobs per second
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			jobID++

			job := Job{
				ID:        jobID,
				Payload:   fmt.Sprintf("job-data-%d", jobID),
				Delay:     time.Duration(rand.Intn(100)) * time.Millisecond,
				LeakSize:  rand.Intn(1024), // Random leak size for testing
				CreatedAt: time.Now(),
			}

			select {
			case app.jobQueue <- job:
				atomic.AddInt64(&app.metrics.PendingJobs, 1)
			default:
				// Queue is full, job dropped
				atomic.AddInt64(&app.metrics.ErrorCount, 1)
			}

		case <-app.shutdown:
			return
		}
	}
}

// processResults handles job completion
func (app *Application) processResults() {
	for {
		select {
		case result := <-app.resultQueue:
			atomic.AddInt64(&app.metrics.PendingJobs, -1)
			atomic.AddInt64(&app.metrics.ProcessedJobs, 1)

			if !result.Success {
				atomic.AddInt64(&app.metrics.ErrorCount, 1)
			}

		case <-app.shutdown:
			return
		}
	}
}

// updateMetrics periodically updates system metrics
func (app *Application) updateMetrics() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			app.metrics.Goroutines = runtime.NumGoroutine()
			app.metrics.MemAllocMB = int(m.Alloc / 1024 / 1024)
			app.metrics.MemSysMB = int(m.Sys / 1024 / 1024)
			app.metrics.GCRuns = int(m.NumGC)
			app.metrics.ActiveWorkers = atomic.LoadInt64(&app.workerPool.activeJobs)
			app.metrics.LeakedGoroutines = atomic.LoadInt64(&app.metrics.LeakedGoroutines)
			app.metrics.BlockedGoroutines = atomic.LoadInt64(&app.metrics.BlockedGoroutines)

		case <-app.shutdown:
			return
		}
	}
}

// simulateMemoryLeak creates a controlled memory leak for testing
func (app *Application) simulateMemoryLeak() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Println("Memory leak simulator started")

	for {
		select {
		case <-ticker.C:
			// Simulate a memory leak - keep growing without bounds
			app.mu.Lock()
			leak := make([]byte, 1024*1024*2) // 2MB each time
			app.memoryLeaks = append(app.memoryLeaks, leak)

			// Intentionally remove the cleanup that was in the original
			// This causes unbounded memory growth
			log.Printf("Memory leak created. Total leaks: %d (approx %dMB)",
				len(app.memoryLeaks), len(app.memoryLeaks)*2)
			app.mu.Unlock()

		case <-app.shutdown:
			return
		}
	}
}

// simulateGoroutineLeak creates goroutines that never exit
func (app *Application) simulateGoroutineLeak() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	log.Println("Goroutine leak simulator started")

	for {
		select {
		case <-ticker.C:
			// Create leaky channels and goroutines
			leakyCh := make(chan int)
			app.leakyChannels = append(app.leakyChannels, leakyCh)

			// Goroutine that blocks forever on channel send
			go func(ch chan int, id int) {
				log.Printf("Starting leaky goroutine %d", id)
				ch <- id // This will block forever since no receiver
				log.Printf("Leaky goroutine %d finished (this should never print)", id)
			}(leakyCh, len(app.leakyChannels))

			// Goroutine that blocks forever on channel receive
			go func(ch chan int, id int) {
				log.Printf("Starting leaky receiver goroutine %d", id)
				<-ch // This will block forever since no sender to this instance
				log.Printf("Leaky receiver %d finished (this should never print)", id)
			}(make(chan int), len(app.leakyChannels))

			atomic.AddInt64(&app.metrics.LeakedGoroutines, 2)
			log.Printf("Created 2 leaky goroutines. Total leaked: %d",
				atomic.LoadInt64(&app.metrics.LeakedGoroutines))

		case <-app.shutdown:
			return
		}
	}
}

// simulateDeadlock creates circular dependency deadlocks
func (app *Application) simulateDeadlock() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Println("Deadlock simulator started")

	for {
		select {
		case <-ticker.C:
			var mutex1, mutex2 sync.Mutex

			log.Println("Creating deadlock scenario...")

			// Goroutine 1: acquires mutex1, then tries mutex2
			go func() {
				log.Println("Goroutine 1: acquiring mutex1")
				mutex1.Lock()
				log.Println("Goroutine 1: acquired mutex1, sleeping...")
				time.Sleep(1 * time.Second)

				log.Println("Goroutine 1: trying to acquire mutex2")
				mutex2.Lock() // Will block
				log.Println("Goroutine 1: acquired both mutexes")
				mutex2.Unlock()
				mutex1.Unlock()
			}()

			// Goroutine 2: acquires mutex2, then tries mutex1
			go func() {
				time.Sleep(500 * time.Millisecond) // Start slightly after goroutine 1
				log.Println("Goroutine 2: acquiring mutex2")
				mutex2.Lock()
				log.Println("Goroutine 2: acquired mutex2, sleeping...")
				time.Sleep(1 * time.Second)

				log.Println("Goroutine 2: trying to acquire mutex1")
				mutex1.Lock() // Will block - DEADLOCK!
				log.Println("Goroutine 2: acquired both mutexes (this should never print)")
				mutex1.Unlock()
				mutex2.Unlock()
			}()

			atomic.AddInt64(&app.metrics.BlockedGoroutines, 2)
			log.Printf("Created deadlock with 2 goroutines. Total blocked: %d",
				atomic.LoadInt64(&app.metrics.BlockedGoroutines))

		case <-app.shutdown:
			return
		}
	}
}

// simulateMutexLeak creates mutexes that are locked but never unlocked
func (app *Application) simulateMutexLeak() {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	log.Println("Mutex leak simulator started")

	for {
		select {
		case <-ticker.C:
			// Create resource with mutex leak
			resource := &MutexLeakResource{
				id:   len(app.mutexLeakResources) + 1,
				data: "initial",
			}
			app.mutexLeakResources = append(app.mutexLeakResources, resource)

			// Goroutine that locks mutex but never unlocks
			go func(r *MutexLeakResource) {
				log.Printf("Resource %d: performing leaky operation", r.id)
				r.LeakyOperation() // Locks but never unlocks
			}(resource)

			// Goroutines that try to access the resource and get blocked
			for i := 0; i < 3; i++ {
				go func(r *MutexLeakResource, accessId int) {
					time.Sleep(2 * time.Second) // Wait for leaky operation to lock
					log.Printf("Resource %d: attempting access %d", r.id, accessId)
					r.ProblematicAccess() // Will block forever
					log.Printf("Resource %d: access %d completed (should never print)", r.id, accessId)
				}(resource, i+1)
			}

			atomic.AddInt64(&app.metrics.BlockedGoroutines, 3)
			log.Printf("Created mutex leak with resource %d. Blocked goroutines: %d",
				resource.id, atomic.LoadInt64(&app.metrics.BlockedGoroutines))

		case <-app.shutdown:
			return
		}
	}
}

// HTTP Handlers

func (app *Application) metricsHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&app.metrics.RequestCount, 1)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.metrics)
}

func (app *Application) healthHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&app.metrics.RequestCount, 1)

	status := map[string]interface{}{
		"status":     "healthy",
		"timestamp":  time.Now().UTC(),
		"goroutines": runtime.NumGoroutine(),
		"version":    "1.0.0",
		"problems_enabled": map[string]bool{
			"memory_leak":    app.enableMemoryLeak,
			"goroutine_leak": app.enableGoroutineLeak,
			"deadlock":       app.enableDeadlock,
			"mutex_leak":     app.enableMutexLeak,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (app *Application) loadHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&app.metrics.RequestCount, 1)

	// Simulate some CPU work
	start := time.Now()
	sum := 0
	for i := 0; i < rand.Intn(1000000); i++ {
		sum += i
	}

	response := map[string]interface{}{
		"message":     "Load endpoint processed",
		"duration_ms": time.Since(start).Milliseconds(),
		"result":      sum,
		"timestamp":   time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// problemsHandler triggers problems on demand
func (app *Application) problemsHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&app.metrics.RequestCount, 1)

	problem := r.URL.Query().Get("type")

	switch problem {
	case "memory":
		// Trigger immediate memory leak
		app.mu.Lock()
		for i := 0; i < 10; i++ {
			leak := make([]byte, 1024*1024*5) // 5MB each
			app.memoryLeaks = append(app.memoryLeaks, leak)
		}
		app.mu.Unlock()
		fmt.Fprintf(w, "Triggered memory leak: 50MB allocated\n")

	case "goroutine":
		// Trigger immediate goroutine leak
		for i := 0; i < 5; i++ {
			deadCh := make(chan int)
			go func(id int) {
				deadCh <- id // Will block forever
			}(i)
		}
		atomic.AddInt64(&app.metrics.LeakedGoroutines, 5)
		fmt.Fprintf(w, "Triggered goroutine leak: 5 goroutines blocked\n")

	case "mutex":
		// Trigger immediate mutex leak
		var mu sync.Mutex
		mu.Lock()
		// Never unlock

		for i := 0; i < 3; i++ {
			go func(id int) {
				mu.Lock() // Will block forever
				defer mu.Unlock()
				log.Printf("Mutex access %d", id)
			}(i)
		}
		atomic.AddInt64(&app.metrics.BlockedGoroutines, 3)
		fmt.Fprintf(w, "Triggered mutex leak: 3 goroutines blocked\n")

	default:
		fmt.Fprintf(w, "Available problems: ?type=memory, ?type=goroutine, ?type=mutex\n")
	}
}

func main() {
	// Create application
	app := NewApplication()

	// Setup HTTP server
	r := mux.NewRouter()

	// Application endpoints
	r.HandleFunc("/health", app.healthHandler).Methods("GET")
	r.HandleFunc("/metrics", app.metricsHandler).Methods("GET")
	r.HandleFunc("/load", app.loadHandler).Methods("GET")
	r.HandleFunc("/problems", app.problemsHandler).Methods("GET")

	// Static load endpoint for testing
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&app.metrics.RequestCount, 1)
		fmt.Fprintf(w, "Go Load Test Application\nGoroutines: %d\nProcessed Jobs: %d\nProblems Enabled: Memory=%v, Goroutine=%v, Deadlock=%v, Mutex=%v\n",
			runtime.NumGoroutine(), atomic.LoadInt64(&app.metrics.ProcessedJobs),
			app.enableMemoryLeak, app.enableGoroutineLeak, app.enableDeadlock, app.enableMutexLeak)
	}).Methods("GET")

	// Main server
	mainServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Debug server (pprof)
	go func() {
		log.Println("Starting debug server on :6060")
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	// Start application
	app.Start()

	// Print usage instructions
	log.Println()
	log.Println("=== Debugging Test Application ===")
	log.Println("Enable problems with environment variables:")
	log.Println("  ENABLE_MEMORY_LEAK=true")
	log.Println("  ENABLE_GOROUTINE_LEAK=true")
	log.Println("  ENABLE_DEADLOCK=true")
	log.Println("  ENABLE_MUTEX_LEAK=true")
	log.Println()
	log.Println("Trigger problems on demand:")
	log.Println("  curl 'http://localhost:8080/problems?type=memory'")
	log.Println("  curl 'http://localhost:8080/problems?type=goroutine'")
	log.Println("  curl 'http://localhost:8080/problems?type=mutex'")
	log.Println()
	log.Println("Debug with pprof:")
	log.Println("  go tool pprof http://localhost:6060/debug/pprof/heap")
	log.Println("  go tool pprof http://localhost:6060/debug/pprof/goroutine")
	log.Println("  curl 'http://localhost:6060/debug/pprof/goroutine?debug=2'")
	log.Println()

	// Graceful shutdown
	go func() {
		log.Printf("Main server listening on :8080")
		log.Printf("Debug server (pprof) listening on :6060")
		if err := mainServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	select {}
}
