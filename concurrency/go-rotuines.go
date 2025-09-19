package main

import (
	"fmt"
	"time"
)

type Person struct {
	Name string
	Age  int
}

func sayHello(p Person) {
	for i := 0; i < 5; i++ {
		fmt.Printf("Hello my name is %s and I am %d years old! (iteration %d)\n", p.Name, p.Age, i+1)
		time.Sleep(100 * time.Millisecond)
	}
}

func GoRoutines() {
	doug := Person{"Doug", 13}
	skeeter := Person{"Skeeter", 13}

	// Sequential
	fmt.Println("=== Sequential ===")
	sayHello(doug)
	sayHello(skeeter)

	// Concurrent execution
	fmt.Println("\n=== Concurrent ===")
	go sayHello(doug)
	go sayHello(skeeter)

	// Wait for goroutines to finish
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Main function ending")
}
