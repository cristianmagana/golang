package main

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"time"

// 	"golang.org/x/time/rate"
// )

// // Mock API response
// type UserData struct {
// 	ID    int
// 	Name  string
// 	Email string
// }

// // Simulates an external API call
// func fetchUserData(ctx context.Context, userID int) (*UserData, error) {

// 	select {
// 	case <-ctx.Done():
// 		return nil, ctx.Err()
// 	//signify network latency
// 	case <-time.After(100 * time.Millisecond):
// 	}
// 	// Simulate network latency
// 	time.Sleep(100 * time.Millisecond)

// 	// Simulate occasional failures
// 	if userID%15 == 0 {
// 		return nil, fmt.Errorf("API error for user %d", userID)
// 	}

// 	return &UserData{
// 		ID:    userID,
// 		Name:  fmt.Sprintf("User%d", userID),
// 		Email: fmt.Sprintf("user%d@example.com", userID),
// 	}, nil
// }

// // TODO: Implement this function
// // It should fetch data for all userIDs concurrently while respecting the rate limit
// func fetchAllUsers(ctx context.Context, userIDs []int, maxRPS int, maxRetries int) (map[int]*UserData, error) {

// 	results := make(map[int]*UserData)
// 	var wg sync.WaitGroup
// 	var mu sync.Mutex

// 	if len(userIDs) == 0 {
// 		return results, nil
// 	}

// 	limiter := rate.NewLimiter(rate.Limit(maxRPS), maxRPS)

// 	for _, userID := range userIDs {
// 		if err := limiter.Wait(ctx); err != nil {
// 			break
// 		}
// 		wg.Add(1)
// 		go func(id int) {
// 			defer wg.Done()
// 			for attempt := 1; attempt <= maxRetries; attempt++ {
// 				select {
// 				case <-ctx.Done():
// 					fmt.Printf("Worker %d context canceled\n", id)
// 					return
// 				default:
// 					data, err := fetchUserData(ctx, id)

// 					if err == nil {
// 						mu.Lock()
// 						results[id] = data
// 						mu.Unlock()
// 						return
// 					}
// 					fmt.Printf("Worker %d failed (attempt %d/%d)\n", id, attempt, maxRetries)

// 					time.Sleep(time.Duration(attempt) * time.Millisecond * 100)
// 				}
// 			}
// 		}(userID)
// 	}

// 	wg.Wait()

// 	return results, ctx.Err()
// }

// func main() {
// 	userIDs := make([]int, 100)
// 	for i := range userIDs {
// 		userIDs[i] = i + 1
// 	}

// 	maxRetries := 3

// 	start := time.Now()
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
// 	defer cancel()

// 	results, err := fetchAllUsers(ctx, userIDs, 10, maxRetries)

// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return
// 	}

// 	for id, data := range results {
// 		fmt.Printf("id %d, data: %v\n", id, data)
// 	}

// 	elapsedTime := time.Since(start)

// 	fmt.Printf("Process took %s\n", elapsedTime)
// 	fmt.Printf("Successfully fetched %d users\n", len(results))
// }
