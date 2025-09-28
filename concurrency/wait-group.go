package main

import (
	"fmt"
	"sync"
	"time"
)

func workerWg(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)
	fmt.Printf("Worker %d done\n", id)
}

func WaitGroup() {
	fmt.Printf("\nStarting waitgroup exercise")
	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1)
		go workerWg(i, &wg)
	}

	wg.Wait()
	fmt.Println("All workers completed")
}
