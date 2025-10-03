package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func producer(ctx context.Context, ch chan<- string) {
	go func() {
		defer close(ch)
		for msg := 0; msg < 20; msg++ {
			select {
			case ch <- fmt.Sprintf("%d\n", msg):
				fmt.Printf("Sending message %d\n", msg)
			case <-ctx.Done():
				fmt.Printf(("Context timeout..."))
				return
			}
		}
	}()
}

func consumer(ctx context.Context, ch <-chan string, wg *sync.WaitGroup, workers int) <-chan struct{} {
	done := make(chan struct{})
	for i := range workers {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for {
				select {
				case msg, ok := <-ch:
					if !ok {
						fmt.Printf("Worker %d: Channel closed, stopping\n", workerId)
						return
					}
					fmt.Printf("Worker %d received message %s", workerId, msg)
				case <-ctx.Done():
					fmt.Printf("Worker %d cancelled due to timeout: %v\n", workerId, ctx.Err())
					return
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	return done
}

func ContextTimeout() {
	ch := make(chan string, 10)
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	producer(ctx, ch)
	done := consumer(ctx, ch, &wg, 10)

	select {
	case <-done:
		fmt.Printf("All workers completed successfully\n")
	case <-ctx.Done():
		fmt.Printf("Processing timed out: %v\n", ctx.Err())
		cancel()
		wg.Wait()
	}

}

// package main

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"time"
// )

// func producer(ctx context.Context, ch chan<- string) {
// 	go func() {
// 		defer close(ch)
// 		for msg := 0; msg < 20; msg++ {
// 			select {
// 			case ch <- fmt.Sprintf("Sending to channel: %d", msg):
// 				// Message sent successfully
// 			case <-ctx.Done():
// 				fmt.Printf("Producer cancelled due to timeout: %v\n", ctx.Err())
// 				return
// 			}
// 		}
// 	}()
// }

// func consumer(ctx context.Context, ch <-chan string, wg *sync.WaitGroup, workers int) <-chan struct{} {
// 	done := make(chan struct{})

// 	for i := 0; i < workers; i++ {
// 		wg.Add(1)
// 		go func(workerId int) {
// 			defer wg.Done()
// 			for {
// 				select {
// 				case msg, ok := <-ch:
// 					if !ok {
// 						return // Channel closed
// 					}
// 					fmt.Printf("Worker %d processed message: %s\n", workerId, msg)
// 				case <-ctx.Done():
// 					fmt.Printf("Worker %d cancelled due to timeout: %v\n", workerId, ctx.Err())
// 					return
// 				}
// 			}
// 		}(i)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(done)
// 	}()

// 	return done
// }

// func ContextTimeout() {
// 	ch := make(chan string, 10)
// 	var wg sync.WaitGroup

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()

// 	// Start producer
// 	producer(ch)

// 	// Start consumers and get done channel
// 	done := consumer(ctx, ch, &wg, 10)

// 	select {
// 	case <-done:
// 		fmt.Println("All workers completed successfully")
// 	case <-ctx.Done():
// 		fmt.Printf("Processing timed out: %v\n", ctx.Err())
// 		// Cancel all workers
// 		cancel()
// 		// Wait for workers to finish gracefully
// 		wg.Wait()
// 	}
// }
