# Go Concurrency Exercises - Progressive Learning

## 🎯 What We're Learning
Go's concurrency model is built around **goroutines** (lightweight threads) and **channels** (communication between goroutines). The philosophy is: *"Don't communicate by sharing memory; share memory by communicating."*

**Key Concepts We'll Master:**
- Goroutines: lightweight concurrent functions
- Channels: typed conduits for data flow
- Channel directions and buffering
- Select statements for multiplexing
- Sync package utilities (WaitGroups, Mutexes)
- Common concurrency patterns

---

## Exercise 1: Your First Goroutine
**Target Understanding:** Basic goroutines and the `go` keyword

### Example:
```go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    for i := 0; i < 3; i++ {
        fmt.Printf("Hello %s! (%d)\n", name, i+1)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Sequential execution
    fmt.Println("=== Sequential ===")
    sayHello("Alice")
    sayHello("Bob")
    
    // Concurrent execution
    fmt.Println("\n=== Concurrent ===")
    go sayHello("Alice")
    go sayHello("Bob")
    
    // Wait for goroutines to finish
    time.Sleep(500 * time.Millisecond)
    fmt.Println("Main function ending")
}
```

### Your Exercise:
Create a program that:
1. Has a function `countDown(name string, seconds int)` that counts down from `seconds` to 1
2. Launches 3 goroutines: "Rocket-1" (5 seconds), "Rocket-2" (3 seconds), "Rocket-3" (7 seconds)
3. Each should print: `"Rocket-1: 5"`, `"Rocket-1: 4"`, etc.
4. Use `time.Sleep(1 * time.Second)` between counts
5. Make sure main waits long enough to see all rockets finish

---

## Exercise 2: Basic Channels
**Target Understanding:** Creating and using unbuffered channels

### Example:
```go
package main

import "fmt"

func sender(ch chan string) {
    ch <- "Hello"
    ch <- "World"
    close(ch) // Signal that no more values will be sent
}

func main() {
    // Create a channel
    messages := make(chan string)
    
    // Start sender goroutine
    go sender(messages)
    
    // Receive from channel
    for msg := range messages {
        fmt.Println("Received:", msg)
    }
}
```

### Your Exercise:
Create a calculator service:
1. Create a `Calculator` struct with `Operation string`, `A int`, `B int`, `Result int`
2. Create channels: `jobs chan Calculator` and `results chan Calculator`
3. Create a worker function that:
   - Receives Calculator from `jobs` channel
   - Performs the operation (add, subtract, multiply)
   - Sends result to `results` channel
4. In main: send 5 different calculations through `jobs`, receive and print 5 results

---

## Exercise 3: Buffered Channels & Channel Directions
**Target Understanding:** Channel buffering and restricting channel directions

### Example:
```go
package main

import (
    "fmt"
    "time"
)

// send-only channel parameter
func producer(ch chan<- int) {
    for i := 1; i <= 5; i++ {
        ch <- i
        fmt.Printf("Sent: %d\n", i)
        time.Sleep(100 * time.Millisecond)
    }
    close(ch)
}

// receive-only channel parameter
func consumer(ch <-chan int) {
    for num := range ch {
        fmt.Printf("Received: %d\n", num)
        time.Sleep(200 * time.Millisecond) // Slower consumer
    }
}

func main() {
    // Buffered channel - can hold 3 values
    numbers := make(chan int, 3)
    
    go producer(numbers)
    consumer(numbers) // Run in main goroutine
}
```

### Your Exercise:
Create a log processing system:
1. Create a `LogEntry` struct: `Timestamp string`, `Level string`, `Message string`
2. Create a buffered channel `logs chan LogEntry` (buffer size 10)
3. Create a `generateLogs` function (send-only channel) that creates 20 log entries
4. Create a `processLogs` function (receive-only channel) that:
   - Filters only "ERROR" and "WARN" level logs
   - Prints them in format: `[TIMESTAMP] LEVEL: MESSAGE`
5. Show how buffering allows producer to get ahead of consumer

---

## Exercise 4: Select Statement
**Target Understanding:** Non-blocking operations and multiplexing channels

### Example:
```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        ch1 <- "Message from channel 1"
    }()
    
    go func() {
        time.Sleep(1 * time.Second)
        ch2 <- "Message from channel 2"
    }()
    
    for i := 0; i < 2; i++ {
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
```

### Your Exercise:
Build a multi-service monitoring system:
1. Create 3 services that send status updates at different intervals:
   - Database service: sends "DB: OK" every 2 seconds
   - API service: sends "API: OK" every 1 second  
   - Cache service: sends "CACHE: OK" every 3 seconds
2. Use select to receive from all services for 10 seconds total
3. Add a timeout case (500ms) that prints "Waiting for updates..."
4. Count and report how many updates received from each service

---

## Exercise 5: WaitGroups
**Target Understanding:** Coordinating multiple goroutines without channels

### Example:
```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // Decrement counter when function returns
    
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Duration(id) * time.Second)
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup
    
    for i := 1; i <= 3; i++ {
        wg.Add(1) // Increment counter
        go worker(i, &wg)
    }
    
    wg.Wait() // Wait for all to finish
    fmt.Println("All workers completed")
}
```

### Your Exercise:
Create a parallel file processor:
1. Create a slice of "filenames": `[]string{"file1.txt", "file2.txt", "file3.txt", "file4.txt", "file5.txt"}`
2. Create a `processFile(filename string, wg *sync.WaitGroup)` function that:
   - Simulates processing with `time.Sleep(time.Duration(len(filename)) * 100 * time.Millisecond)`
   - Prints processing start and completion
3. Process all files concurrently using goroutines and WaitGroup
4. Measure and print total execution time vs. sequential processing

---

## Exercise 6: Mutexes and Shared State
**Target Understanding:** Protecting shared data from race conditions

### Example:
```go
package main

import (
    "fmt"
    "sync"
)

type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

func main() {
    counter := &Counter{}
    var wg sync.WaitGroup
    
    // 100 goroutines incrementing 100 times each
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 100; j++ {
                counter.Increment()
            }
        }()
    }
    
    wg.Wait()
    fmt.Printf("Final counter value: %d\n", counter.Value()) // Should be 10000
}
```

### Your Exercise:
Build a thread-safe bank account system:
1. Create a `BankAccount` struct with `balance int` and `mutex sync.RWMutex`
2. Implement methods:
   - `Deposit(amount int)` - adds to balance
   - `Withdraw(amount int) bool` - subtracts if sufficient funds, returns success
   - `Balance() int` - returns current balance (use read lock)
3. Create simulation with:
   - 5 depositor goroutines (each deposits $100, 10 times)
   - 3 withdrawer goroutines (each withdraws $50, 20 times)
   - 1 balance checker (prints balance every 100ms)
4. Run for 2 seconds, show final balance

---

## Exercise 7: Worker Pool Pattern
**Target Understanding:** Managing a pool of workers processing jobs

### Example:
```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Job struct {
    ID     int
    Data   string
    Result string
}

func worker(id int, jobs <-chan Job, results chan<- Job, wg *sync.WaitGroup) {
    defer wg.Done()
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job.ID)
        time.Sleep(100 * time.Millisecond) // Simulate work
        job.Result = fmt.Sprintf("Processed by worker %d: %s", id, job.Data)
        results <- job
    }
}

func main() {
    jobs := make(chan Job, 10)
    results := make(chan Job, 10)
    
    // Start 3 workers
    var wg sync.WaitGroup
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(i, jobs, results, &wg)
    }
    
    // Send jobs
    for i := 1; i <= 9; i++ {
        jobs <- Job{ID: i, Data: fmt.Sprintf("data-%d", i)}
    }
    close(jobs)
    
    // Close results when all workers done
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    for result := range results {
        fmt.Printf("Result: %s\n", result.Result)
    }
}
```

### Your Exercise:
Build an image processing service:
1. Create an `Image` struct: `ID int`, `Name string`, `Size int`, `Status string`
2. Create processing functions that take 200-500ms each:
   - `resize(img *Image)` - sets Status to "resized"
   - `compress(img *Image)` - sets Status to "compressed" 
   - `watermark(img *Image)` - sets Status to "watermarked"
3. Create a pipeline: resize → compress → watermark
4. Use 2 workers for each stage, process 15 images
5. Track and print processing time for entire batch

---

## Exercise 8: Fan-in/Fan-out Pattern
**Target Understanding:** Distributing work and combining results

### Your Exercise:
Create a distributed prime number finder:
1. **Fan-out**: Split number range 1-1000 into chunks, send to multiple workers
2. **Workers**: Each checks if numbers in their chunk are prime
3. **Fan-in**: Collect all prime numbers from workers into single channel
4. Use 5 workers, each processing 200 numbers
5. Print all primes found and total count
6. Compare execution time vs single-threaded approach

---

## 🏆 Final Challenge: Web Scraper with Concurrency
Combine everything you've learned:

1. **Input**: List of 10 URLs (use httpbin.org endpoints for testing)
2. **Rate limiting**: Maximum 3 concurrent requests
3. **Timeout**: Each request times out after 2 seconds  
4. **Results**: Collect response status codes and response times
5. **Error handling**: Handle timeouts and HTTP errors gracefully
6. **Reporting**: Print summary of successful/failed requests

**Concepts to use:**
- Worker pool for rate limiting
- Channels for job distribution and result collection
- Select with timeouts
- WaitGroups for coordination
- Mutexes for shared counters

---

## 💡 Learning Tips:

1. **Run each exercise multiple times** - concurrency can behave differently each run
2. **Add logging** - Use `fmt.Printf` liberally to see execution order
3. **Experiment with timing** - Change sleep durations to see different behaviors
4. **Use race detector** - Run with `go run -race main.go` to catch race conditions
5. **Start simple** - Get basic version working before adding complexity

## 🔧 Common Patterns You'll Learn:

- **Pipeline**: Data flows through stages (Exercise 7)
- **Fan-out/Fan-in**: Distribute work, collect results (Exercise 8)
- **Worker Pool**: Fixed number of workers processing jobs (Exercise 7)
- **Publish/Subscribe**: Multiple consumers of same data
- **Rate Limiting**: Control resource usage (Final Challenge)

Work through these exercises in order - each builds on concepts from the previous ones!