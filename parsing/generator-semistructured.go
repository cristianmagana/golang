package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var (
	logLevels = []string{"INFO", "WARN", "ERROR", "DEBUG"}
	services  = []string{
		"auth-service",
		"api-gateway",
		"database-service",
		"cache-service",
		"payment-service",
		"notification-service",
		"user-service",
		"order-service",
	}

	messages = map[string][]string{
		"ERROR": {
			"Connection timeout",
			"Query failed",
			"Authentication failed",
			"Connection lost",
			"Token expired",
			"Timeout",
			"Database unreachable",
			"Memory allocation failed",
			"Request rejected",
			"Service unavailable",
		},
		"WARN": {
			"Slow response",
			"High memory usage",
			"Retry attempt",
			"Deprecated API used",
			"Rate limit approaching",
			"Cache miss",
			"Queue filling up",
		},
		"INFO": {
			"Request processed",
			"User logged in",
			"Cache updated",
			"Configuration loaded",
			"Health check passed",
			"Job completed",
			"Connection established",
		},
		"DEBUG": {
			"Entering function",
			"Variable value",
			"Step completed",
			"Query executed",
		},
	}

	// Weight distribution for log levels (INFO most common, ERROR least)
	levelWeights = map[string]int{
		"INFO":  60,
		"DEBUG": 25,
		"WARN":  10,
		"ERROR": 5,
	}
)

// GenerateLogFile creates a log file with specified number of entries
func GenerateLogFile(filename string, numEntries int) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Start time
	startTime := time.Date(2024, 9, 29, 10, 0, 0, 0, time.UTC)
	currentTime := startTime

	fmt.Printf("Generating %d log entries...\n", numEntries)
	progressInterval := numEntries / 10

	for i := 0; i < numEntries; i++ {
		if i > 0 && i%progressInterval == 0 {
			fmt.Printf("Progress: %d%%\n", (i*100)/numEntries)
		}

		// Generate timestamp (increment by 1-10 seconds randomly)
		increment := rand.Intn(10) + 1
		currentTime = currentTime.Add(time.Duration(increment) * time.Second)
		timestamp := currentTime.Format("2006-01-02 15:04:05")

		// Select log level based on weights
		level := getWeightedLevel()

		// Select random service
		service := services[rand.Intn(len(services))]

		// Select appropriate message for level
		messageList := messages[level]
		message := messageList[rand.Intn(len(messageList))]

		// Write log line
		logLine := fmt.Sprintf("%s %s %s %s\n", timestamp, level, service, message)
		_, err := writer.WriteString(logLine)
		if err != nil {
			return fmt.Errorf("failed to write log line: %v", err)
		}
	}

	fmt.Println("Progress: 100%")
	fmt.Printf("Successfully generated %d log entries in %s\n", numEntries, filename)

	return nil
}

// getWeightedLevel returns a log level based on weights
func getWeightedLevel() string {
	totalWeight := 0
	for _, weight := range levelWeights {
		totalWeight += weight
	}

	r := rand.Intn(totalWeight)
	cumulative := 0

	for level, weight := range levelWeights {
		cumulative += weight
		if r < cumulative {
			return level
		}
	}

	return "INFO" // fallback
}

func main10() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate different sized log files
	testSizes := []struct {
		filename string
		entries  int
	}{
		{"small_logs.txt", 1_000},
		{"medium_logs.txt", 100_000},
		{"large_logs.txt", 1_000_000},
	}

	for _, test := range testSizes {
		fmt.Printf("\n=== Generating %s ===\n", test.filename)
		start := time.Now()

		err := GenerateLogFile(test.filename, test.entries)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		elapsed := time.Since(start)
		fmt.Printf("Generation time: %v\n", elapsed)

		// Get file size
		fileInfo, _ := os.Stat(test.filename)
		sizeMB := float64(fileInfo.Size()) / (1024 * 1024)
		fmt.Printf("File size: %.2f MB\n", sizeMB)
	}

	fmt.Println("\n=== Summary ===")
	fmt.Println("Generated files:")
	fmt.Println("  small_logs.txt   - 1,000 entries")
	fmt.Println("  medium_logs.txt  - 100,000 entries")
	fmt.Println("  large_logs.txt   - 1,000,000 entries")
	fmt.Println("\nNow you can test your parser with:")
	fmt.Println("  stats, _ := ProcessLogsFromFileConcurrent(\"large_logs.txt\", 4)")
}
