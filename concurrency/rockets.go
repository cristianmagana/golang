package main

import (
	"fmt"
	"time"
)

func countDown(name string, seconds int) {
	for i := seconds; i > 0; i-- {
		fmt.Printf("%s: %d\n", name, i)
		time.Sleep(1 * time.Second)
	}
}

func Rockets() {
	fmt.Println()
	fmt.Println("Ground control this is major Cristian")

	go countDown("Rocket-1", 5)
	go countDown("Rocket-2", 3)
	go countDown("Rocket-3", 7)

	time.Sleep(8 * time.Second)
	fmt.Println("All rockets have launched!")
}
