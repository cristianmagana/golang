package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main2() {
	// Generate different sized files for testing
	sizes := map[string]int{
		"small.csv":  1000,    // 1K rows
		"medium.csv": 100000,  // 100K rows
		"large.csv":  1000000, // 1M rows
		"huge.csv":   5000000, // 5M rows
	}

	services := []string{
		"user-service", "order-service", "payment-service",
		"auth-service", "inventory-service", "notification-service",
		"analytics-service", "reporting-service", "billing-service",
	}

	for filename, rowCount := range sizes {
		fmt.Printf("Generating %s with %d rows...\n", filename, rowCount)
		start := time.Now()

		file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}

		// Write header
		file.WriteString("timestamp,service,response_time_ms,memory_mb,cpu_percent,disk_read_kb,disk_write_kb,network_kb\n")

		// Generate data
		baseTime := time.Date(2025, 9, 25, 10, 0, 0, 0, time.UTC)

		for i := 0; i < rowCount; i++ {
			// Add some malformed data occasionally (1% of rows)
			if rand.Intn(100) == 0 {
				file.WriteString("malformed,data,here,missing,fields\n")
				continue
			}

			timestamp := baseTime.Add(time.Duration(i) * time.Second).Format(time.RFC3339)
			service := services[rand.Intn(len(services))]
			responseTime := rand.Float64()*100 + 10 // 10-110ms
			memory := rand.Intn(2048) + 256         // 256-2304 MB
			cpu := rand.Intn(100)                   // 0-100%
			diskRead := rand.Intn(4096) + 512       // 512-4608 KB
			diskWrite := rand.Intn(6144) + 1024     // 1024-7168 KB
			network := rand.Intn(500) + 50          // 50-550 KB

			line := fmt.Sprintf("%s,%s,%.1f,%d,%d,%d,%d,%d\n",
				timestamp, service, responseTime, memory, cpu, diskRead, diskWrite, network)

			file.WriteString(line)

			// Progress indicator for large files
			if rowCount > 100000 && i%100000 == 0 {
				fmt.Printf("  Progress: %d/%d rows\n", i, rowCount)
			}
		}

		file.Close()
		elapsed := time.Since(start)
		fmt.Printf("  Completed %s in %v\n\n", filename, elapsed)
	}

	fmt.Println("All CSV files generated!")
	fmt.Println("\nFile sizes:")
	for filename := range sizes {
		if stat, err := os.Stat(filename); err == nil {
			fmt.Printf("  %s: %.1f MB\n", filename, float64(stat.Size())/1024/1024)
		}
	}
}
