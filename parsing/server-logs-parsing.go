package main

/*
2024-09-29 10:00:07 INFO order-service Health check passed
2024-09-29 10:00:10 WARN user-service Queue filling up
2024-09-29 10:00:14 INFO auth-service Connection established
*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

type LogEntry struct {
	Timestamp   string `json:"timestamp"`
	Level       string `json:"level"`
	Application string `json:"application"`
	Message     string `json:"message"`
}

func processLogEntry(line string) (*LogEntry, error) {
	parts := strings.SplitN(line, " ", 4)
	if len(parts) > 4 {
		return nil, fmt.Errorf("error reading the log line")
	}

	timeStamp := parts[0] + " " + parts[1]

	level := parts[2]

	serviceParts := strings.SplitN(parts[3], " ", 2)
	if len(serviceParts) != 2 {
		return nil, fmt.Errorf("error reading the log line")
	}

	application := serviceParts[0]
	message := serviceParts[1]

	return &LogEntry{
		Timestamp:   timeStamp,
		Level:       level,
		Application: application,
		Message:     message,
	}, nil
}

func lineParser(lines chan<- string, errChan chan<- error, path string) {
	go func() {
		defer close(lines)

		file, err := os.Open(path)
		if err != nil {
			errChan <- err
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan()

		for scanner.Scan() {
			lines <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			errChan <- err
		}
		close(errChan)
	}()
}

func lineWorker(lines <-chan string, logEntry chan<- *LogEntry, workers int, errChan chan<- error) {
	var wg sync.WaitGroup

	for i := range workers {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()

			for line := range lines {
				entry, err := processLogEntry(line)
				if err != nil {
					errChan <- err
					return
				}
				logEntry <- entry
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(logEntry)
	}()
}

func main() {
	lines := make(chan string, 100)
	logEntry := make(chan *LogEntry, 100)
	errChan := make(chan error, 1)

	path := "./medium_logs.txt"

	lineParser(lines, errChan, path)
	lineWorker(lines, logEntry, 10, errChan)

	// Handle both log entries and errors
	go func() {
		for err := range errChan {
			fmt.Printf("Error: %v\n", err)
		}
	}()

	fmt.Print("[")
	first := true
	for result := range logEntry {
		if !first {
			fmt.Print(",")
		}
		jsonData, err := json.Marshal(result)
		if err != nil {
			fmt.Printf("error marshalling %v", err)
			continue
		}
		fmt.Print(string(jsonData))
		first = false
	}
	fmt.Println("]")

}
