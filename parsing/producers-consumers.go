package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// This is going to put data onto the buffer stream
func fileParser(lines chan<- string, errc chan<- error, path string) {

	go func() {

		defer close(lines)

		file, err := os.Open(path)
		if err != nil {
			errc <- err
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan()
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			errc <- err
		}
		close(errc)
	}()
}

func parserOutput(results chan<- string, lines <-chan string) {

	var wg sync.WaitGroup
	workers := 10

	for i := range workers {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()

			for line := range lines {
				results <- line
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()
}

func main4() {
	lines := make(chan string, 1000)
	results := make(chan string, 100)
	errc := make(chan error, 1)

	fileParser(lines, errc, "./medium.csv")
	parserOutput(results, lines)

	for result := range results {
		fmt.Println(result)
	}
}
