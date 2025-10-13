package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Notification represents a message to send
type Notification struct {
	ID       string
	Type     string // "email", "sms", "push"
	Message  string
	Attempts int
}

// Simple rate limiter using token bucket
type RateLimiter struct {
	tokens   int
	max      int
	rate     int // tokens per second
	mu       sync.Mutex
	lastTime time.Time
}

func NewRateLimiter(max, rate int) *RateLimiter {
	return &RateLimiter{
		tokens:   max,
		max:      max,
		rate:     rate,
		lastTime: time.Now(),
	}
}

func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(r.lastTime).Seconds()
	tokensToAdd := int(elapsed * float64(r.rate))

	if tokensToAdd > 0 {
		r.tokens = min(r.max, r.tokens+tokensToAdd)
		r.lastTime = now
	}

	// Check if we have tokens
	if r.tokens > 0 {
		r.tokens--
		return true
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Dispatcher manages workers and queues
type Dispatcher struct {
	jobs       chan *Notification
	retryQueue chan *Notification
	limiters   map[string]*RateLimiter
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewDispatcher(workers int) *Dispatcher {
	ctx, cancel := context.WithCancel(context.Background())

	d := &Dispatcher{
		jobs:       make(chan *Notification, 100),
		retryQueue: make(chan *Notification, 50),
		limiters: map[string]*RateLimiter{
			"email": NewRateLimiter(100, 10),
			"sms":   NewRateLimiter(50, 5),
			"push":  NewRateLimiter(200, 20),
		},
		ctx:    ctx,
		cancel: cancel,
	}

	// Start workers
	for i := 0; i < workers; i++ {
		d.wg.Add(1)
		go d.worker(i)
	}

	// Start retry handler
	d.wg.Add(1)
	go d.retryHandler()

	return d
}

func (d *Dispatcher) worker(id int) {
	defer d.wg.Done()

	for {
		select {
		case <-d.ctx.Done():
			log.Printf("Worker %d stopping", id)
			return
		case notif := <-d.jobs:
			d.process(id, notif)
		}
	}
}

func (d *Dispatcher) process(workerID int, notif *Notification) {
	// Check rate limit
	limiter := d.limiters[notif.Type]
	if !limiter.Allow() {
		log.Printf("Worker %d: Rate limit hit for %s, queuing retry", workerID, notif.Type)
		d.scheduleRetry(notif)
		return
	}

	// Simulate sending (with 20% failure rate)
	time.Sleep(50 * time.Millisecond)
	failed := time.Now().UnixNano()%5 == 0

	notif.Attempts++

	if failed {
		log.Printf("Worker %d: Failed to send %s (attempt %d)", workerID, notif.ID, notif.Attempts)
		if notif.Attempts < 3 {
			d.scheduleRetry(notif)
		} else {
			log.Printf("Worker %d: Max retries reached for %s", workerID, notif.ID)
		}
		return
	}

	log.Printf("Worker %d: ✓ Sent %s via %s", workerID, notif.ID, notif.Type)
}

func (d *Dispatcher) scheduleRetry(notif *Notification) {
	go func() {
		// Exponential backoff: 1s, 2s, 4s
		backoff := time.Duration(1<<uint(notif.Attempts)) * time.Second

		select {
		case <-time.After(backoff):
			d.retryQueue <- notif
		case <-d.ctx.Done():
		}
	}()
}

func (d *Dispatcher) retryHandler() {
	defer d.wg.Done()

	for {
		select {
		case <-d.ctx.Done():
			log.Println("Retry handler stopping")
			return
		case notif := <-d.retryQueue:
			d.jobs <- notif
		}
	}
}

func (d *Dispatcher) Submit(notif *Notification) error {
	select {
	case d.jobs <- notif:
		return nil
	case <-d.ctx.Done():
		return fmt.Errorf("dispatcher shutting down")
	default:
		return fmt.Errorf("queue full")
	}
}

func (d *Dispatcher) Shutdown(timeout time.Duration) error {
	log.Println("Shutting down...")
	d.cancel()

	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Shutdown complete")
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("shutdown timeout")
	}
}

func Notifications() {
	dispatcher := NewDispatcher(5) // 5 workers

	// Submit test notifications
	types := []string{"email", "sms", "push"}
	for i := 0; i < 30; i++ {
		notif := &Notification{
			ID:      fmt.Sprintf("notif-%d", i),
			Type:    types[i%3],
			Message: fmt.Sprintf("Test message %d", i),
		}

		if err := dispatcher.Submit(notif); err != nil {
			log.Printf("Failed to submit: %v", err)
		}
	}

	// Let it run
	time.Sleep(8 * time.Second)

	// Graceful shutdown
	if err := dispatcher.Shutdown(3 * time.Second); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
}
