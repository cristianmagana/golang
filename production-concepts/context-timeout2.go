package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func producerContext(ctx context.Context, ch chan<- string) {
	go func() {
		defer close(ch)
		for msg := 0; msg < 100; msg++ {

			select {
			case ch <- fmt.Sprintf("%d\n", msg):
				fmt.Printf("Sending message %d\n", msg)
			case <-ctx.Done():
				fmt.Printf("Context timeout exceeded\n")
				return
			}
		}

	}()
}

func consumerContext(ctx context.Context, ch <-chan string, wg *sync.WaitGroup, workers int) <-chan struct{} {

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

func ContextTimeout2() {
	ch := make(chan string)

	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	producerContext(ctx, ch)

	done := consumerContext(ctx, ch, &wg, 10)

	select {
	case <-done:
		fmt.Println("All workers completed successfully")
	case <-ctx.Done():
		fmt.Printf("Context timeout exceeded\n")
		cancel()
		wg.Wait()
	}
}
