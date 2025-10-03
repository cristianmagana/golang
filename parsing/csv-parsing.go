package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sync"
	"time"
)

// Add some expensive processing to see concurrency benefits
func expensiveProcessing(line string) string {
	// Simulate CPU-intensive work
	var result float64
	for i := 0; i < 5000; i++ {
		result += math.Sin(float64(i)) * math.Cos(float64(i))
	}

	// Return processed line with some indicator
	return fmt.Sprintf("[PROCESSED %.2f] %s", result, line)
}

func main3() {
	lines := make(chan string, 100)
	results := make(chan string, 1000)
	errorc := make(chan error, 1)

	start := time.Now()

	// file reading concurrency
	go func() {
		defer close(lines)

		file, err := os.Open("./large.csv")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan()
		for scanner.Scan() {
			lines <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			errorc <- err
		}
		close(errorc)
	}()

	// Worker pools

	var wg sync.WaitGroup
	workers := 50
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()

			for line := range lines {
				// does expensive processing
				processed := expensiveProcessing(line)
				results <- processed
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Printf("%s\n", result)
	}

	select {
	case err := <-errorc:
		if err != nil {
			fmt.Println("Error: ", err)
		}
	default:
	}
	fmt.Printf("Took: %d", time.Since(start).Milliseconds())

}
