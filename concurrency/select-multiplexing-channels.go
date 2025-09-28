package main

import (
	"fmt"
	"time"
)

func Select() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		//time.Sleep(2 * time.Second)
		ch1 <- "Message from channel 1"
	}()

	go func() {
		//time.Sleep(1 * time.Second)
		ch2 <- "Message from channel 2"
	}()

	for i := 0; i < 100; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println(msg1)
		case msg2 := <-ch2:
			fmt.Println(msg2)
		case <-time.After(3 * time.Second):
			fmt.Println("Timeout!")
		}
	}
}
