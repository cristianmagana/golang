package main

import (
	"fmt"
)

type Calculator struct {
	Operation string
	A         int
	B         int
	Result    int
}

func worker(jobs chan Calculator, results chan Calculator) {
	for job := range jobs {
		switch job.Operation {
		case "add":
			job.Result = job.A + job.B
		case "subtract":
			job.Result = job.A - job.B
		case "multiply":
			job.Result = job.A * job.B
		default:
			job.Result = 0
		}
		results <- job
	}
}

func getOperationSymbol(operation string) string {
	switch operation {
	case "add":
		return "+"
	case "subtract":
		return "-"
	case "multiply":
		return "*"
	default:
		return "?"
	}
}

func CalculatorConcurrency() {
	fmt.Println()
	fmt.Println("Calculator")

	jobs := make(chan Calculator)
	results := make(chan Calculator)

	go worker(jobs, results)

	go func() {
		jobs <- Calculator{Operation: "add", A: 10, B: 5}
		jobs <- Calculator{Operation: "subtract", A: 20, B: 8}
		jobs <- Calculator{Operation: "multiply", A: 7, B: 6}
		jobs <- Calculator{Operation: "add", A: 15, B: 25}
		jobs <- Calculator{Operation: "multiply", A: 0, B: 1}
		close(jobs)
	}()

	// Receive and print 5 results
	for range 5 {
		result := <-results
		fmt.Printf("Operation: %s, %d %s %d = %d\n",
			result.Operation, result.A, getOperationSymbol(result.Operation), result.B, result.Result)
	}

}
