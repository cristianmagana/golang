package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Mock API response
type UserData struct {
	ID    int
	Name  string
	Email string
}

// Simulates an external API call
func fetchUserData(ctx context.Context, userID int) (*UserData, error) {

	select {
	case <-time.After(100 * time.Millisecond):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Simulate occasional failures
	if userID%15 == 0 {
		return nil, fmt.Errorf("API error for user %d", userID)
	}

	return &UserData{
		ID:    userID,
		Name:  fmt.Sprintf("User%d", userID),
		Email: fmt.Sprintf("user%d@example.com", userID),
	}, nil
}

func fetchUserDataWithRetry(ctx context.Context, userID int, maxRetries int) (*UserData, error) {

	var lastError error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		data, err := fetchUserData(ctx, userID)
		if err == nil {
			return data, nil
		}

		lastError = err

		if attempt < maxRetries-1 {
			backoff := time.Duration(attempt*100) * time.Millisecond
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, lastError)
}

// TODO: Implement this function
// It should fetch data for all userIDs concurrently while respecting the rate limit
func fetchAllUsers(ctx context.Context, userIDs []int, maxRPS int) (map[int]*UserData, error) {

	if len(userIDs) == 0 {
		return make(map[int]*UserData), nil
	}

	results := make(map[int]*UserData)
	var mu sync.Mutex
	var wg sync.WaitGroup

	limiter := rate.NewLimiter(rate.Limit(maxRPS), maxRPS)

	for _, userID := range userIDs {
		if err := limiter.Wait(ctx); err != nil {
			break
		}

		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			data, err := fetchUserDataWithRetry(ctx, id, 3)
			if err == nil {
				mu.Lock()
				results[id] = data
				mu.Unlock()
			}
		}(userID)
	}
	wg.Wait()
	return results, ctx.Err()
}

func main3() {
	userIDs := make([]int, 100)
	for i := range userIDs {

		userIDs[i] = i + 1
		//fmt.Printf("%d\n", userIDs[i])
	}

	maxRPS := 10
	estimatedTime := time.Duration(len(userIDs)/maxRPS) * time.Second
	timeout := estimatedTime * 2 // ✅ Add buffer

	ctx, cancel := context.WithTimeout(context.Background(), timeout*30)
	defer cancel()

	start := time.Now()
	results, err := fetchAllUsers(ctx, userIDs, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	elapsed := time.Since(start)

	count := 0
	for _, data := range results {
		b, _ := json.Marshal(data)
		fmt.Println(string(b))
		count++
	}

	fmt.Printf("Successfully fetched %d/%d users in %v\n", len(results), len(userIDs), elapsed)
}
