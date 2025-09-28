package main

import (
	"fmt"
	"time"
)

func main() {

	start := time.Now()
	//GoRoutines()
	//Rockets()
	//CalculatorConcurrency()

	fmt.Println("=== Your Original (Always Ordered) ===")
	// SenderReceiverChannels()

	// Uncomment these to see out-of-order behavior:
	// MultipleWorkers()
	// MultipleGoroutinesRacing()
	// ModifiedSenderReceiver()

	// PingPong()

	// BufferedChannelLogs()

	// Select()

	// WaitGroup()

	// WaitGroupServices()

	// MutexSharedState()

	MutexSharedBankService()

	fmt.Printf("\nProcess took: %s", time.Since(start))

}
