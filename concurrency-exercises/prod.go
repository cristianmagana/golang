package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// SHARED STATE - Protected by RWMutex
// ============================================================================

type Analytics struct {
	mu              sync.RWMutex
	requestCount    map[string]int64   // endpoint -> count
	errorCount      map[string]int64   // error type -> count
	avgResponseTime map[string]float64 // endpoint -> avg ms
}

func NewAnalytics() *Analytics {
	return &Analytics{
		requestCount:    make(map[string]int64),
		errorCount:      make(map[string]int64),
		avgResponseTime: make(map[string]float64),
	}
}

func (a *Analytics) RecordRequest(endpoint string, duration time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.requestCount[endpoint]++

	// Calculate running average
	count := float64(a.requestCount[endpoint])
	currentAvg := a.avgResponseTime[endpoint]
	newAvg := (currentAvg*(count-1) + float64(duration.Milliseconds())) / count
	a.avgResponseTime[endpoint] = newAvg
}

func (a *Analytics) RecordError(errorType string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.errorCount[errorType]++
}

func (a *Analytics) GetStats() map[string]interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Copy data while holding read lock
	stats := map[string]interface{}{
		"request_count":     copyMap(a.requestCount),
		"error_count":       copyMap(a.errorCount),
		"avg_response_time": copyMapFloat(a.avgResponseTime),
	}

	return stats
}

func copyMap(m map[string]int64) map[string]int64 {
	copy := make(map[string]int64, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return copy
}

func copyMapFloat(m map[string]float64) map[string]float64 {
	copy := make(map[string]float64, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return copy
}

// ============================================================================
// CACHE - Protected by RWMutex with TTL
// ============================================================================

type CacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

type Cache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
}

func NewCache() *Cache {
	c := &Cache{
		entries: make(map[string]*CacheEntry),
	}

	// Background cleanup goroutine
	go c.cleanupExpired()

	return c
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists || time.Now().After(entry.expiresAt) {
		return nil, false
	}

	return entry.value, true
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = &CacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) cleanupExpired() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.entries {
			if now.After(entry.expiresAt) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

// ============================================================================
// TASK QUEUE - Producer-Consumer with Channels
// ============================================================================

type Task struct {
	ID       string
	Type     string
	Payload  map[string]interface{}
	ResultCh chan TaskResult
}

type TaskResult struct {
	TaskID string
	Result interface{}
	Error  error
}

type TaskProcessor struct {
	taskQueue      chan Task
	numWorkers     int
	wg             sync.WaitGroup
	analytics      *Analytics
	tasksProcessed int64
}

func NewTaskProcessor(queueSize, numWorkers int, analytics *Analytics) *TaskProcessor {
	return &TaskProcessor{
		taskQueue:  make(chan Task, queueSize),
		numWorkers: numWorkers,
		analytics:  analytics,
	}
}

func (tp *TaskProcessor) Start(ctx context.Context) {
	// Start worker pool
	for i := 0; i < tp.numWorkers; i++ {
		tp.wg.Add(1)
		go tp.worker(ctx, i)
	}
}

func (tp *TaskProcessor) worker(ctx context.Context, id int) {
	defer tp.wg.Done()

	log.Printf("Worker %d started", id)

	for {
		select {
		case task := <-tp.taskQueue:
			result := tp.processTask(task)

			// Send result back
			select {
			case task.ResultCh <- result:
			case <-time.After(1 * time.Second):
				log.Printf("Worker %d: timeout sending result for task %s", id, task.ID)
			}

			atomic.AddInt64(&tp.tasksProcessed, 1)

		case <-ctx.Done():
			log.Printf("Worker %d shutting down", id)
			return
		}
	}
}

func (tp *TaskProcessor) processTask(task Task) TaskResult {
	// Simulate different task types with varying processing times
	var result interface{}
	var err error

	switch task.Type {
	case "compute":
		time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)
		result = map[string]interface{}{
			"computation": "completed",
			"value":       rand.Intn(1000),
		}
	case "database":
		time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
		result = map[string]interface{}{
			"records": rand.Intn(100),
		}
	case "external_api":
		time.Sleep(time.Duration(200+rand.Intn(300)) * time.Millisecond)
		if rand.Float32() < 0.1 {
			err = fmt.Errorf("external API timeout")
			tp.analytics.RecordError("external_api_timeout")
		} else {
			result = map[string]interface{}{
				"status": "success",
			}
		}
	default:
		err = fmt.Errorf("unknown task type: %s", task.Type)
		tp.analytics.RecordError("unknown_task_type")
	}

	return TaskResult{
		TaskID: task.ID,
		Result: result,
		Error:  err,
	}
}

func (tp *TaskProcessor) SubmitTask(task Task) error {
	select {
	case tp.taskQueue <- task:
		return nil
	case <-time.After(100 * time.Millisecond):
		return fmt.Errorf("task queue full")
	}
}

func (tp *TaskProcessor) Stop() {
	close(tp.taskQueue)
	tp.wg.Wait()
}

func (tp *TaskProcessor) GetProcessedCount() int64 {
	return atomic.LoadInt64(&tp.tasksProcessed)
}

// ============================================================================
// RATE LIMITER - Semaphore (Buffered Channel)
// ============================================================================

type RateLimiter struct {
	sem           chan struct{}
	maxConcurrent int
	activeCount   int64
}

func NewRateLimiter(maxConcurrent int) *RateLimiter {
	return &RateLimiter{
		sem:           make(chan struct{}, maxConcurrent),
		maxConcurrent: maxConcurrent,
	}
}

func (rl *RateLimiter) Acquire() {
	rl.sem <- struct{}{}
	atomic.AddInt64(&rl.activeCount, 1)
}

func (rl *RateLimiter) Release() {
	<-rl.sem
	atomic.AddInt64(&rl.activeCount, -1)
}

func (rl *RateLimiter) GetActiveCount() int64 {
	return atomic.LoadInt64(&rl.activeCount)
}

// ============================================================================
// HTTP SERVER
// ============================================================================

type Server struct {
	analytics     *Analytics
	cache         *Cache
	taskProcessor *TaskProcessor
	rateLimiter   *RateLimiter
	reqCounter    int64
}

func NewServer() *Server {
	analytics := NewAnalytics()

	return &Server{
		analytics:     analytics,
		cache:         NewCache(),
		taskProcessor: NewTaskProcessor(1000, 50, analytics), // 50 workers, 1000 queue size
		rateLimiter:   NewRateLimiter(100),                   // Max 100 concurrent requests
	}
}

func (s *Server) Start(ctx context.Context) {
	// Start background workers
	s.taskProcessor.Start(ctx)

	// HTTP handlers
	http.HandleFunc("/api/process", s.handleProcess)
	http.HandleFunc("/api/cached-data", s.handleCachedData)
	http.HandleFunc("/api/stats", s.handleStats)
	http.HandleFunc("/health", s.handleHealth)

	// Start server
	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	s.taskProcessor.Stop()
	log.Println("Server stopped")
}

// ============================================================================
// HANDLER: Process Task (Rate Limited + Producer-Consumer)
// ============================================================================

func (s *Server) handleProcess(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	reqID := atomic.AddInt64(&s.reqCounter, 1)

	// Rate limiting with semaphore
	s.rateLimiter.Acquire()
	defer s.rateLimiter.Release()

	// Track analytics
	defer func() {
		duration := time.Since(start)
		s.analytics.RecordRequest("/api/process", duration)
	}()

	// Parse request
	var req struct {
		TaskType string                 `json:"task_type"`
		Payload  map[string]interface{} `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.analytics.RecordError("invalid_json")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create task
	task := Task{
		ID:       fmt.Sprintf("task-%d", reqID),
		Type:     req.TaskType,
		Payload:  req.Payload,
		ResultCh: make(chan TaskResult, 1),
	}

	// Submit to worker pool (producer)
	if err := s.taskProcessor.SubmitTask(task); err != nil {
		s.analytics.RecordError("queue_full")
		http.Error(w, "Service busy", http.StatusServiceUnavailable)
		return
	}

	// Wait for result (consumer)
	select {
	case result := <-task.ResultCh:
		if result.Error != nil {
			s.analytics.RecordError("task_failed")
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"task_id":         result.TaskID,
			"result":          result.Result,
			"active_requests": s.rateLimiter.GetActiveCount(),
		})

	case <-time.After(5 * time.Second):
		s.analytics.RecordError("task_timeout")
		http.Error(w, "Request timeout", http.StatusGatewayTimeout)
	}
}

// ============================================================================
// HANDLER: Cached Data (RWMutex for Cache)
// ============================================================================

func (s *Server) handleCachedData(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	s.rateLimiter.Acquire()
	defer s.rateLimiter.Release()

	defer func() {
		s.analytics.RecordRequest("/api/cached-data", time.Since(start))
	}()

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}

	// Try cache first (RWMutex read)
	if value, found := s.cache.Get(key); found {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "HIT")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"key":    key,
			"value":  value,
			"cached": true,
		})
		return
	}

	// Cache miss - simulate expensive operation
	time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)

	value := map[string]interface{}{
		"data":      fmt.Sprintf("computed_value_%d", time.Now().Unix()),
		"timestamp": time.Now().Unix(),
	}

	// Store in cache (RWMutex write)
	s.cache.Set(key, value, 30*time.Second)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"key":    key,
		"value":  value,
		"cached": false,
	})
}

// ============================================================================
// HANDLER: Stats (RWMutex Read)
// ============================================================================

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := s.analytics.GetStats()

	stats["tasks_processed"] = s.taskProcessor.GetProcessedCount()
	stats["active_requests"] = s.rateLimiter.GetActiveCount()
	stats["max_concurrent"] = s.rateLimiter.maxConcurrent
	stats["worker_count"] = s.taskProcessor.numWorkers

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// ============================================================================
// HANDLER: Health Check
// ============================================================================

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}

// ============================================================================
// LOAD TESTER
// ============================================================================

func runLoadTest(baseURL string, duration time.Duration) {
	log.Println("Starting load test...")

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var wg sync.WaitGroup
	requestCount := int64(0)
	errorCount := int64(0)

	// Simulate 200 concurrent clients
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			client := &http.Client{Timeout: 10 * time.Second}

			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Random endpoint
					var resp *http.Response
					var err error

					switch rand.Intn(3) {
					case 0:
						// Process task
						//body := `{"task_type":"compute","payload":{"value":123}}`
						resp, err = client.Post(baseURL+"/api/process", "application/json",
							http.NoBody)
					case 1:
						// Cached data
						key := fmt.Sprintf("key_%d", rand.Intn(10))
						resp, err = client.Get(baseURL + "/api/cached-data?key=" + key)
					case 2:
						// Stats
						resp, err = client.Get(baseURL + "/api/stats")
					}

					atomic.AddInt64(&requestCount, 1)

					if err != nil {
						atomic.AddInt64(&errorCount, 1)
						time.Sleep(10 * time.Millisecond)
						continue
					}

					resp.Body.Close()

					// Small delay between requests
					time.Sleep(time.Duration(10+rand.Intn(20)) * time.Millisecond)
				}
			}
		}(i)
	}

	// Stats reporter
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for range ticker.C {
			req := atomic.LoadInt64(&requestCount)
			err := atomic.LoadInt64(&errorCount)
			log.Printf("Load test - Requests: %d, Errors: %d, Error Rate: %.2f%%",
				req, err, float64(err)/float64(req)*100)
		}
	}()

	wg.Wait()
	ticker.Stop()

	total := atomic.LoadInt64(&requestCount)
	errors := atomic.LoadInt64(&errorCount)
	log.Printf("Load test complete - Total: %d, Errors: %d, Success Rate: %.2f%%",
		total, errors, float64(total-errors)/float64(total)*100)
}

// ============================================================================
// MAIN
// ============================================================================

func prod() {
	server := NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start server
	go server.Start(ctx)

	// Wait for server to start
	time.Sleep(1 * time.Second)

	// Run load test
	runLoadTest("http://localhost:8080", 30*time.Second)

	// Let server run a bit more
	time.Sleep(5 * time.Second)

	// Shutdown
	cancel()
	time.Sleep(2 * time.Second)
}
