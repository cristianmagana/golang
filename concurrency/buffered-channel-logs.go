package main

import (
	"fmt"
	"time"
)

type LogEntry struct {
	Timestamp string
	Level     string
	Message   string
}

func generateLogs(log chan<- LogEntry) {
	for i := range 200 {
		// Method 1: Using fmt.Sprintf (most common and flexible)
		timestamp := time.Now().Format(time.RFC3339)
		message := fmt.Sprintf("Hello from %d", i)

		l := LogEntry{timestamp, "INFO", message}
		if i%2 > 0 {
			l.Level = "ERROR"
		}
		fmt.Printf("Sent %s [Timestamp: %s, Level: %s, Message: %s]\n", time.Now().Format(time.RFC3339), l.Timestamp, l.Level, l.Message)
		log <- l
		time.Sleep(2000 * time.Millisecond)
	}
	close(log)
}

func processLogs(log <-chan LogEntry) {
	for l := range log {
		if l.Level == "ERROR" {
			fmt.Printf("Received %s [Timestamp: %s, Level: %s, Message: %s]\n", time.Now().Format(time.RFC3339), l.Timestamp, l.Level, l.Message)
			time.Sleep(10000 * time.Millisecond)
		}
	}
}

func BufferedChannelLogs() {
	fmt.Println("\nStarting logs processor")

	logs := make(chan LogEntry, 10)

	go generateLogs(logs)
	processLogs(logs)

}
