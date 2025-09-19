# 🚀 Go Learning Journey: From Primitives to Production

A comprehensive learning path for mastering Go (Golang) through hands-on exercises covering core language features, concurrency patterns, and production debugging techniques.

## 🎯 Learning Objectives

This repository is designed to take you from Go beginner to confident practitioner through three progressive learning modules:

1. **🔧 [Primitives](./primatives/)** - Master Go's fundamental building blocks
2. **⚡ [Concurrency](./concurrency/)** - Harness Go's powerful concurrency model  
3. **🔍 [Performance Profiling](./pprof/)** - Debug and optimize production applications

## 📚 Learning Path

### Phase 1: Master the Fundamentals
**📁 [Go Primitives →](./primatives/readme.md)**

Build a solid foundation with Go's core data types and language constructs:

- **Variables & Data Types** - strings, integers, booleans, runes
- **Collections** - arrays, slices, maps with practical operations
- **Control Flow** - if/else, loops, switch statements
- **Functions & Methods** - parameters, returns, receivers
- **Pointers** - memory addresses and dereferencing
- **Structs** - custom types and method receivers

*Perfect for: Newcomers to Go or those needing a refresher on fundamentals*

### Phase 2: Embrace Concurrency
**📁 [Go Concurrency →](./concurrency/readme.md)**

Learn Go's signature feature through progressive exercises:

- **Goroutines** - lightweight concurrent functions
- **Channels** - safe communication between goroutines
- **Select Statements** - multiplexing and timeouts
- **Sync Package** - WaitGroups, Mutexes, coordination
- **Concurrency Patterns** - worker pools, fan-in/fan-out, pipelines
- **Real-world Applications** - web scrapers, data processors

*Perfect for: Developers ready to leverage Go's concurrency superpowers*

### Phase 3: Production-Ready Debugging
**📁 [Performance Profiling with pprof →](./pprof/readme.md)**

Master production debugging and optimization:

- **CPU Profiling** - identify performance bottlenecks
- **Memory Profiling** - detect and fix memory leaks
- **Goroutine Analysis** - debug deadlocks and leaks
- **Kubernetes Integration** - profile containerized applications
- **Performance Optimization** - systematic troubleshooting workflows

*Perfect for: Developers deploying Go applications to production*

## 🛠️ Getting Started

### Prerequisites
- Go 1.19+ installed ([download here](https://golang.org/dl/))
- Basic programming experience (any language)
- Terminal/command line familiarity

### Quick Start
```bash
# Clone the repository
git clone <your-repo-url>
cd golang

# Start with primitives
cd primatives
go run main.go

# Move to concurrency when ready
cd ../concurrency  
go run main.go

# Finally, explore profiling
cd ../pprof
go run main.go
```

## 📖 How to Use This Repository

### 🎯 **Recommended Learning Sequence**

1. **Start with [Primitives](./primatives/readme.md)** - Even experienced developers should review Go's unique features
2. **Progress to [Concurrency](./concurrency/readme.md)** - Work through exercises sequentially
3. **Apply [Profiling](./pprof/readme.md)** - Practice with the provided sample application

### 🏃‍♂️ **Exercise Approach**

- **Read** the README for each section thoroughly
- **Code along** with examples before attempting exercises
- **Experiment** - modify examples to see different behaviors  
- **Build** complete solutions, don't just read code
- **Test** your understanding with the provided challenges

### 🔄 **Iterative Learning**

- Come back to earlier sections after completing later ones
- Use concurrent patterns in your primitive exercises
- Apply profiling to your concurrency solutions
- Real mastery comes from connecting concepts across all three areas

## 📋 Progress Tracking

- [ ] **Primitives Completed** - Comfortable with Go basics
- [ ] **First Goroutine** - Written concurrent code
- [ ] **Channel Communication** - Goroutines talking to each other
- [ ] **Production Debugging** - Used pprof to find real issues
- [ ] **Complete Application** - Built something using all three areas

## 🎯 What You'll Build

By the end of this journey, you'll have:

- **Solid Go Foundation** - Confident with syntax and idioms
- **Concurrency Skills** - Able to write safe, efficient concurrent programs
- **Debugging Expertise** - Capable of troubleshooting production issues
- **Real Projects** - Calculator services, web scrapers, monitoring systems
- **Production Readiness** - Ready to build and deploy Go applications

## 💡 Learning Tips

### 🏗️ **Build, Don't Just Read**
- Type out every example
- Modify code to see what breaks
- Create your own variations

### 🧪 **Experiment Freely**
- Go is safe to experiment with
- Use `go run -race` to catch concurrency issues
- Break things on purpose to understand them

### 🔄 **Connect the Dots**
- See how primitives enable concurrency
- Use profiling to understand your concurrent code
- Apply learnings across all three areas

### 📚 **Reference, Don't Memorize**
- Keep the READMEs open while coding
- Bookmark useful patterns
- Focus on understanding over memorization

## 🚀 Next Steps

After completing these exercises:

1. **Build a Real Project** - Combine all three areas in a web service
2. **Contribute to Open Source** - Your Go skills are now valuable
3. **Explore Advanced Topics** - Context, reflection, code generation
4. **Share Your Knowledge** - Teach others what you've learned

---

## 📂 Repository Structure

```
golang/
├── README.md                          # This file - start here!
├── primatives/                        # Phase 1: Go fundamentals
│   ├── readme.md                      # Complete primitives guide
│   ├── main.go                        # Runnable examples
│   └── *.go                           # Individual concept files
├── concurrency/                       # Phase 2: Concurrent programming  
│   ├── readme.md                      # Progressive concurrency exercises
│   ├── main.go                        # Basic examples
│   └── *.go                           # Specific pattern implementations
├── pprof/                             # Phase 3: Production debugging
│   ├── readme.md                      # Complete profiling guide
│   ├── main.go                        # Sample application to profile
│   ├── Dockerfile                     # Container setup
│   └── k8s-deployment.yaml            # Kubernetes integration
└── into/                              # Legacy examples (reference only)
```

Start your Go journey now with **[Go Primitives →](./primatives/readme.md)**

Happy coding! 🎉
