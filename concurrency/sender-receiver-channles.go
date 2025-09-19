package main

import (
	"fmt"
	"time"
)

func producer(ch chan<- int) {
	for i := range 100 {
		ch <- i
		fmt.Printf("Sent %d\n", i)
		time.Sleep(10 * time.Millisecond)
	}
	close(ch)
}

func consumer(ch <-chan int) {
	for i := range ch {
		fmt.Printf("Received %d\n", i)
		time.Sleep(2000 * time.Millisecond)
	}
}

func SenderReceiverChannels() {
	fmt.Println("\nStarting Sender and Receiver channels")

	numbers := make(chan int, 5)

	go producer(numbers)
	consumer(numbers)
}
