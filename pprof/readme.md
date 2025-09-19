# Go Application Troubleshooting with pprof

A comprehensive guide to using pprof for debugging CPU performance, memory leaks, goroutine issues, deadlocks, and concurrency problems in Go applications.

## Prerequisites

```bash
# Ensure your application has pprof enabled
import _ "net/http/pprof"

# Start debug server
go func() {
    log.Println(http.ListenAndServe(":6060", nil))
}()

# Port forward if running in Kubernetes
kubectl port-forward service/your-app-debug 6060:6060
kubectl port-forward service/your-app-service 8080:8080
```

## CPU Profiling

### Collecting CPU Profile

```bash
# Collect 30-second CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

### Understanding CPU Profile Output

```
Duration: 30.01s, Total samples = 360ms (1.20%)
```

- **Duration**: How long profiling ran (30 seconds)
- **Total samples**: Actual CPU time captured (360ms)
- **Percentage**: CPU utilization (1.20% = low CPU usage, not CPU-bound)

### Analyzing CPU Profile

```bash
# Inside pprof interactive mode
(pprof) top
Showing nodes accounting for 360ms, 100% of 360ms total
      flat  flat%   sum%        cum   cum%
     300ms 83.33% 83.33%      300ms 83.33%  main.processJob
      20ms  5.56% 88.89%       20ms  5.56%  runtime.futex
      10ms  2.78% 97.22%       10ms  2.78%  runtime.lock2
```

#### Column Definitions:
- **flat**: Time spent directly in this function
- **flat%**: Percentage of total CPU time in this function  
- **sum%**: Cumulative percentage up to this line
- **cum**: Cumulative time including functions this calls
- **cum%**: Percentage of total including called functions

#### Key Commands:

```bash
(pprof) top 10                    # Top 10 CPU consumers
(pprof) top -cum                  # Sort by cumulative time
(pprof) list function_name        # Source code with CPU annotations
(pprof) peek function_name        # Show callers and callees
(pprof) traces                    # Call stack traces
(pprof) web                       # Visual call graph (requires Graphviz)
```

### CPU Profile Red Flags

**Performance Issues:**
- Functions with unexpectedly high `flat%` or `cum%`
- Many samples in `runtime.*` functions (GC pressure, lock contention)
- High CPU usage in goroutine scheduling (`runtime.findRunnable`)

**Lock Contention:**
```bash
(pprof) top
# Look for:
runtime.semacquire1     # Mutex contention
runtime.lock2          # Runtime lock contention
sync.(*Mutex).Lock     # Application mutex contention
```

## Heap Profiling

### Different Heap Profile Types

```bash
# Current memory usage (default)
go tool pprof http://localhost:6060/debug/pprof/heap

# Total allocations over lifetime  
go tool pprof -alloc_space http://localhost:6060/debug/pprof/heap

# Total objects allocated over lifetime
go tool pprof -alloc_objects http://localhost:6060/debug/pprof/heap

# Current objects in memory
go tool pprof -inuse_objects http://localhost:6060/debug/pprof/heap
```

### Analyzing Heap Profiles

```bash
# Inside pprof interactive mode
(pprof) top
Showing nodes accounting for 34210.05kB, 100% of 34210.05kB total
      flat  flat%   sum%        cum   cum%
33159.67kB 96.93% 96.93% 33159.67kB 96.93%  main.simulateMemoryLeak
  536.37kB  1.57% 98.50%   536.37kB  1.57%  main.NewApplication
     514kB  1.50%   100%      514kB  1.50%  bufio.NewWriterSize
```

#### Key Commands:

```bash
(pprof) top                       # Biggest memory consumers
(pprof) top -cum                  # Sort by cumulative usage
(pprof) list function_name        # Source code with memory annotations
(pprof) traces                    # Allocation stack traces
(pprof) sample_index alloc_space  # Switch to total allocations
(pprof) sample_index inuse_space  # Switch back to current usage
(pprof) sample_index alloc_objects # Switch to object count
```

### Sample Index Analysis

#### `inuse_space` (Default)
- Shows **current memory usage**
- Best for identifying **active memory leaks**
- Use when memory keeps growing

#### `alloc_space`
- Shows **total allocations** over application lifetime
- Best for identifying **allocation hot spots**
- Use when you suspect high allocation rate

#### `inuse_objects`
- Shows **current number of objects** in memory
- Best for identifying **object leaks**
- Use when you suspect specific object types accumulating

#### `alloc_objects`
- Shows **total objects allocated** over lifetime
- Best for identifying **allocation patterns**
- Use to find frequently allocated objects

### Memory Leak Detection

```bash
# Capture baseline
go tool pprof -raw -output=heap1.pprof http://localhost:6060/debug/pprof/heap

# Wait 5-10 minutes during normal operation
sleep 600

# Capture comparison profile  
go tool pprof -raw -output=heap2.pprof http://localhost:6060/debug/pprof/heap

# Compare to see what grew
go tool pprof -base=heap1.pprof heap2.pprof
(pprof) top    # Shows memory that increased between captures
```

### Memory Leak Patterns

**Unbounded Growth:**
```
main.cacheHandler        50MB  # Map/cache growing without cleanup
main.processLogs         25MB  # Slice accumulating without bounds  
main.websocketHandler    15MB  # Connections not being cleaned up
```

**High Allocation Rate:**
```bash
(pprof) sample_index alloc_space
(pprof) top
# Look for functions with very high allocation rates
# that might benefit from object pooling
```

## Goroutine Profiling

### Collecting Goroutine Profile

```bash
# Interactive analysis
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Quick goroutine count
curl "http://localhost:6060/debug/pprof/goroutine?debug=0" | wc -l

# Detailed goroutine dump with stack traces
curl "http://localhost:6060/debug/pprof/goroutine?debug=2" | head -100

# Goroutine states summary
curl "http://localhost:6060/debug/pprof/goroutine?debug=1"
```

### Analyzing Goroutine Profiles

```bash
# Inside pprof interactive mode
(pprof) top
Showing nodes accounting for 15 goroutines, 100% of 15 total
      flat  flat%   sum%        cum   cum%
         9  60.00% 60.00%          9  60.00%  main.(*WorkerPool).worker
         2  13.33% 73.33%          2  13.33%  main.(*Application).processResults
         1   6.67% 80.00%          1   6.67%  main.(*Application).generateJobs
```

#### Key Commands:

```bash
(pprof) top              # Functions where most goroutines are blocked
(pprof) traces           # Stack traces of all goroutines
(pprof) list worker      # Source code showing where workers block
(pprof) peek processJob  # Function call relationships
```

### Understanding Goroutine States

```bash
# Check goroutine states
curl "http://localhost:6060/debug/pprof/goroutine?debug=1" | \
  grep -E "^goroutine [0-9]+" | \
  sed 's/.*\[\(.*\)\].*/\1/' | \
  sort | uniq -c | sort -nr

# Example output:
#   9 select           # Normal - waiting in select statement  
#   3 chan receive     # Waiting for channel data
#   2 IO wait         # Waiting for network/file I/O
#   1 semacquire      # Blocked on mutex
```

#### Common Goroutine States:

- **`select`**: Waiting in select statement (normal)
- **`chan receive`**: Blocked waiting for channel data
- **`chan send`**: Blocked trying to send to full channel
- **`semacquire`**: Blocked on mutex/semaphore
- **`IO wait`**: Blocked on network/file operations (normal)
- **`sleep`**: In time.Sleep()
- **`runnable`**: Ready to run (many = CPU bottleneck)

## Identifying Specific Problems

### Goroutine Leaks

**Symptoms:**
- Steadily increasing goroutine count
- Many goroutines in same blocking state
- Goroutines that never complete

**Detection:**
```bash
# Monitor goroutine count over time
watch -n 5 'curl -s "http://localhost:6060/debug/pprof/goroutine?debug=0" | wc -l'

# Look for blocked goroutines
curl "http://localhost:6060/debug/pprof/goroutine?debug=2" | \
  grep -A5 -B5 "chan send\|chan receive"
```

**Leak Patterns in Traces:**
```bash
(pprof) traces
# Look for repeated identical stack traces:
         5   runtime.gopark
             runtime.chansend1
             main.leakyFunction
# ^ 5 goroutines stuck sending to same channel = leak
```

### Channel Deadlocks

**Blocked Channel Send:**
```bash
# Detailed stack trace shows:
goroutine 123 [chan send]:
main.problematicSender(0xc000054060)
    /app/main.go:45 +0x85
```

**Blocked Channel Receive:**
```bash
goroutine 124 [chan receive]:  
main.waitingReceiver(0xc000054070)
    /app/main.go:60 +0x92
```

**Detection Pattern:**
- Many goroutines blocked on `chan send` with no corresponding `chan receive`
- Or vice versa - receivers with no senders

### Mutex Deadlocks

**Semacquire Pattern:**
```bash
# Look for goroutines blocked on mutex acquisition
curl "http://localhost:6060/debug/pprof/goroutine?debug=2" | \
  grep -A10 -B5 semacquire

# Example deadlock pattern:
goroutine 45 [semacquire]:
sync.(*Mutex).Lock(0xc00008e000)
main.functionA()
    /app/main.go:100 +0x123

goroutine 46 [semacquire]:  
sync.(*Mutex).Lock(0xc00008e020)
main.functionB()
    /app/main.go:200 +0x456
```

**In pprof traces:**
```bash
(pprof) traces
# Multiple goroutines blocked on semacquire:
         3   runtime.gopark
             runtime.semacquire1
             sync.(*Mutex).lockSlow
             sync.(*Mutex).Lock
             main.problematicFunction
```

### Mutex Never Released

**Detection:**
- High count of goroutines blocked on same mutex location
- Goroutines that never progress from `semacquire` state
- Same function appearing repeatedly in blocked traces

**Example Pattern:**
```bash
(pprof) top
# Shows many goroutines blocked in same function:
        15  100%   100%         15   100%  main.(*Resource).Access
```

## Reading Traces Effectively

### CPU Traces
```bash
(pprof) traces
# Shows call stack samples with frequency:
File: main
Type: cpu
-----------+-------------------------------------------------------
        10   main.expensiveFunction
             main.processData
             main.handleRequest
             main.main
-----------+-------------------------------------------------------
         5   runtime.mallocgc
             main.allocateMemory
             main.processData
```

**Interpretation:**
- First number (10, 5) = sample count
- Higher counts = more CPU time spent
- Stack shows call chain leading to expensive operation

### Heap Traces  
```bash
(pprof) traces
# Shows allocation stack traces:
-----------+-------------------------------------------------------
     50MB   main.cacheData
             main.processRequest
             main.handleHTTP
-----------+-------------------------------------------------------
     25MB   main.logMessage
             main.debugFunction
```

**Interpretation:**
- Size (50MB, 25MB) = memory allocated
- Stack shows where allocations originated
- Helps identify allocation hot spots

### Goroutine Traces
```bash
(pprof) traces  
# Shows where goroutines are blocked:
-----------+-------------------------------------------------------
         9   runtime.gopark
             runtime.selectgo  
             main.(*WorkerPool).worker
-----------+-------------------------------------------------------
         3   runtime.gopark
             runtime.chansend1
             main.leakyFunction
```

**Interpretation:**
- Number (9, 3) = goroutine count in this state
- Stack shows blocking operation
- `runtime.gopark` = goroutine parked (blocked)
- `runtime.selectgo` = blocked in select
- `runtime.chansend1` = blocked sending to channel

## Complete Debugging Workflow

### 1. Establish Baseline
```bash
# Capture normal state
curl http://localhost:8080/metrics | jq
go tool pprof -raw -output=baseline_heap.pprof http://localhost:6060/debug/pprof/heap
go tool pprof -raw -output=baseline_goroutine.pprof http://localhost:6060/debug/pprof/goroutine
```

### 2. Generate Load
```bash
# Stress test your application
for i in {1..100}; do curl -s http://localhost:8080/load > /dev/null & done
```

### 3. Monitor for Issues
```bash
# Watch for growing metrics
watch -n 10 'curl -s http://localhost:8080/metrics | jq "{goroutines, mem_alloc_mb}"'

# Monitor goroutine states
watch -n 5 'curl -s "http://localhost:6060/debug/pprof/goroutine?debug=1" | grep -E "^goroutine [0-9]+" | sed "s/.*\[\(.*\)\].*/\1/" | sort | uniq -c'
```

### 4. Capture Problem State
```bash
# After issues manifest
go tool pprof -raw -output=problem_heap.pprof http://localhost:6060/debug/pprof/heap
go tool pprof -raw -output=problem_goroutine.pprof http://localhost:6060/debug/pprof/goroutine
```

### 5. Analyze Differences
```bash
# Compare profiles
go tool pprof -base=baseline_heap.pprof problem_heap.pprof
go tool pprof -base=baseline_goroutine.pprof problem_goroutine.pprof
```

## Red Flags Summary

### CPU Issues
- Functions consuming unexpected CPU time
- High `runtime.*` function usage
- Lock contention patterns

### Memory Issues  
- Steady growth in `inuse_space` without plateau
- Functions dominating memory allocation
- High `alloc_space` to `inuse_space` ratio

### Goroutine Issues
- Steadily increasing goroutine count
- Many goroutines in same blocking state  
- Goroutines blocked on channel operations without counterparts
- Multiple goroutines blocked on same mutex

### Deadlock Indicators
- Goroutines stuck in `semacquire` that never progress
- Circular dependency patterns in mutex acquisition
- Channel operations with no corresponding send/receive operations

Use this guide systematically to identify and resolve performance issues, memory leaks, and concurrency problems in your Go applications.