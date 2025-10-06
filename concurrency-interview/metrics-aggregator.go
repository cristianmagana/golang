package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// ServiceMetrics represents metrics from a single service
type ServiceMetrics struct {
	ServiceID   string
	CPUUsage    float64
	MemoryUsage float64
	RequestRate int
}

// fetchServiceMetrics simulates fetching metrics from an external service
// DO NOT MODIFY THIS FUNCTION
func fetchServiceMetrics(serviceID string) (*ServiceMetrics, error) {
	// Simulate network latency (varies between 50-150ms)
	latency := 50 + (time.Now().UnixNano()%100)*int64(time.Millisecond)
	time.Sleep(time.Duration(latency))

	// Simulate failures for services ending in "5" (service-5, service-15, etc.)
	if serviceID[len(serviceID)-1] == '5' {
		return nil, fmt.Errorf("service %s unavailable", serviceID)
	}

	// Return mock metrics
	return &ServiceMetrics{
		ServiceID:   serviceID,
		CPUUsage:    float64(time.Now().UnixNano()%100) / 100.0,
		MemoryUsage: float64(time.Now().UnixNano()%80) / 100.0,
		RequestRate: int(time.Now().UnixNano() % 1000),
	}, nil
}

// TODO: Implement this function
// It should fetch metrics from all services concurrently while respecting the rate limit
func aggregateMetrics(ctx context.Context, serviceIDs []string, maxRPS int) (map[string]*ServiceMetrics, error) {
	// Your implementation here

	result := make(map[string]*ServiceMetrics)
	var mu sync.Mutex
	var wg sync.WaitGroup

	limiter := rate.NewLimiter(rate.Limit(maxRPS), maxRPS)
	semaphore := make(chan struct{}, maxRPS)

	for _, serviceId := range serviceIDs {
		if err := limiter.Wait(ctx); err != nil {
			break
		}

		select {
		case semaphore <- struct{}{}:

		case <-ctx.Done():
			goto done
		}
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			data, err := fetchServiceMetrics(id)
			if err == nil {
				mu.Lock()
				result[id] = data
				mu.Unlock()
			}

		}(serviceId)
	}
done:
	wg.Wait()
	return result, ctx.Err()
}

func MetricsAggregator() {
	// Generate 50 service IDs
	serviceIDs := make([]string, 50)
	for i := range serviceIDs {
		serviceIDs[i] = fmt.Sprintf("service-%d", i+1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	start := time.Now()
	results, err := aggregateMetrics(ctx, serviceIDs, 5)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Successfully fetched metrics from %d/%d services in %v\n",
		len(results), len(serviceIDs), elapsed)
}
