package main

import (
	"fmt"
	"sync"
	"time"
)

// DataRecord represents raw input data
type DataRecord struct {
	ID    int
	Value float64
}

// ProcessedData represents data after processing
type ProcessedData struct {
	ID             int
	OriginalValue  float64
	ProcessedValue float64
	ProcessorID    int
}

// AggregatedResult represents final aggregated output
type AggregatedResult struct {
	TotalRecords int
	Sum          float64
	Average      float64
}

// fetchData simulates reading data from a source (database, file, API, etc.)
// DO NOT MODIFY THIS FUNCTION
func fetchData(recordID int) (*DataRecord, error) {
	time.Sleep(50 * time.Millisecond) // Simulate I/O

	// Simulate 10% failure rate
	if recordID%10 == 0 {
		return nil, fmt.Errorf("failed to fetch record %d", recordID)
	}

	return &DataRecord{
		ID:    recordID,
		Value: float64(recordID * 10),
	}, nil
}

// processRecord simulates expensive processing (encryption, transformation, etc.)
// DO NOT MODIFY THIS FUNCTION
func processRecord(record *DataRecord, processorID int) *ProcessedData {
	time.Sleep(100 * time.Millisecond) // Simulate heavy computation

	return &ProcessedData{
		ID:             record.ID,
		OriginalValue:  record.Value,
		ProcessedValue: record.Value * 2.5, // Some transformation
		ProcessorID:    processorID,
	}
}

// TODO: Implement this function (FAN-OUT)
// Read records and fan them out to multiple channels for parallel processing
// Should:
// - Fetch all records with IDs from 1 to numRecords
// - Send successfully fetched records to the output channel
// - Close the channel when done
// - Handle fetch errors gracefully (skip failed records)
func dataSource(numRecords int, maxConcurrency int) <-chan *DataRecord {
	// Your implementation here
	out := make(chan *DataRecord)
	var wg sync.WaitGroup
	sema := make(chan struct{}, maxConcurrency)

	for i := 1; i <= numRecords; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sema <- struct{}{}
			defer func() { <-sema }()
			rec, err := fetchData(id)
			if err == nil {
				out <- rec
			}

		}(i)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// TODO: Implement this function (PROCESSING STAGE)
// Process records from input channel and send to output channel
// Should:
// - Read from input channel until closed
// - Process each record
// - Send processed data to output channel
// - Close output channel when input is closed
func processor(processorID int, input <-chan *DataRecord) <-chan *ProcessedData {
	out := make(chan *ProcessedData)

	go func() {
		defer close(out)
		for record := range input {
			out <- processRecord(record, processorID)
		}
	}()

	return out
}

// TODO: Implement this function (FAN-IN)
// Merge multiple processed data channels into a single channel
// Should:
// - Read from all input channels concurrently
// - Send all data to single output channel
// - Close output channel when ALL inputs are closed
func merge(channels ...<-chan *ProcessedData) <-chan *ProcessedData {
	// Your implementation here
	out := make(chan *ProcessedData)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan *ProcessedData) {
			defer wg.Done()
			for data := range c {
				out <- data
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// TODO: Implement this function (AGGREGATOR)
// Aggregate all processed data into final result
// Should:
// - Read from input channel until closed
// - Calculate total records, sum, and average
func aggregate(input <-chan *ProcessedData) *AggregatedResult {
	// Your implementation here
	var total int
	var sum float64
	for data := range input {
		total++
		sum += data.ProcessedValue
	}

	avg := 0.0
	if total > 0 {
		avg = sum / float64(total)
	}

	return &AggregatedResult{
		TotalRecords: total,
		Sum:          sum,
		Average:      avg,
	}
}

func main_() {
	numRecords := 100
	numProcessors := 5
	maxConcurrency := 1

	start := time.Now()

	// TODO: Build the pipeline
	// 1. Create data source (fan-out point)
	source := dataSource(numRecords, maxConcurrency)
	// 2. Create multiple processors
	var processors []<-chan *ProcessedData
	for i := 1; i <= numProcessors; i++ {
		processors = append(processors, processor(i, source))
	}

	// 3. Merge processor outputs (fan-in point)
	merged := merge(processors...)

	// 4. Aggregate final results
	result := aggregate(merged)

	// Your pipeline construction here

	// Print results
	elapsed := time.Since(start)
	fmt.Printf("Processed %d records in %v\n", result.TotalRecords, elapsed)
	fmt.Printf("Sum: %.2f, Avg: %.2f\n", result.Sum, result.Average)
}
