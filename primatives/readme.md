# Go Primitives Quick Reference

A quick, practical guide to Go's core data types and operations.

## 📝 Basic Data Types

### Variable Declaration
```go
// Explicit declaration
var name string = "Alice"
var age int = 30

// Short declaration (preferred)
name := "Alice"
age := 30
```

### Core Types
```go
// String
message := "Hello, World!"

// Integers  
count := 42

// Float
price := 19.99

// Boolean
isActive := true

// Rune (single character/Unicode code point)
grade := 'A'
```

### Type Formatting
```go
fmt.Printf("String: %s\n", "hello")     // String: hello
fmt.Printf("Integer: %d\n", 42)         // Integer: 42
fmt.Printf("Float: %f\n", 3.14)         // Float: 3.140000
fmt.Printf("Boolean: %t\n", true)       // Boolean: true
fmt.Printf("Rune: %c\n", 'A')          // Rune: A
```

---

## 🔤 String Operations

```go
import "strings"

sentence := "  Hello, World! Welcome to Go Programming  "

// Common operations
strings.TrimSpace(sentence)              // Remove whitespace
strings.ToUpper(sentence)                // UPPERCASE
strings.Replace(sentence, "World", "Go", -1)  // Replace all
strings.Split(sentence, " ")             // Split into slice
strings.Contains(sentence, "Go")         // Check if contains
len(sentence)                            // String length
```

---

## 🔢 Math Operations

```go
import "math"

a, b := 10, 3

// Basic arithmetic
fmt.Println(a + b)    // 13 (addition)
fmt.Println(a - b)    // 7  (subtraction)
fmt.Println(a * b)    // 30 (multiplication)
fmt.Println(a % b)    // 1  (modulo)

// Division differences
fmt.Println(a / b)              // 3 (integer division)
fmt.Println(float64(a) / float64(b))  // 3.333... (float division)

// Math package functions
math.Sqrt(64)        // 8.0
math.Pow(2, 3)       // 8.0
math.Ceil(3.2)       // 4.0
math.Floor(3.8)      // 3.0
```

---

## 📋 Arrays & Slices

### Arrays (Fixed Size)
```go
// Declaration
numbers := [5]int{1, 2, 3, 4, 5}
fmt.Println(numbers)  // [1 2 3 4 5]
```

### Slices (Dynamic)
```go
// From array
slice := numbers[:3]  // [1 2 3]

// Direct creation
colors := []string{"red", "blue", "green"}

// Append elements
colors = append(colors, "yellow")

// Iteration
for i, color := range colors {
    fmt.Printf("Index %d: %s\n", i, color)
}

// Properties
len(colors)  // Length
cap(colors)  // Capacity
```

---

## 🗺️ Maps (Key-Value)

```go
// Create map
grades := make(map[string]string)

// Add elements
grades["Alice"] = "A"
grades["Bob"] = "B"
grades["Charlie"] = "C"

// Update
grades["Charlie"] = "A"

// Check existence
value, exists := grades["Alice"]
if exists {
    fmt.Println("Alice's grade:", value)
}

// Delete
delete(grades, "Bob")

// Iterate
for name, grade := range grades {
    fmt.Printf("%s: %s\n", name, grade)
}
```

---

## 🔄 Control Flow

### If-Else
```go
guess := 25
target := 42

if guess < target {
    fmt.Println("Too low!")
} else if guess > target {
    fmt.Println("Too high!")
} else {
    fmt.Println("Correct!")
}
```

### For Loops
```go
// Standard loop
for i := 1; i <= 10; i++ {
    fmt.Println(i)
}

// Range loop (slices/maps)
for index, value := range slice {
    fmt.Printf("Index: %d, Value: %v\n", index, value)
}

// While-style loop
count := 0
for count < 5 {
    fmt.Println(count)
    count++
}
```

### Switch Statements
```go
number := 42

// Expression switch
switch number % 2 {
case 0:
    fmt.Println("Even")
default:
    fmt.Println("Odd")
}

// Boolean switch
switch {
case number < 10:
    fmt.Println("Small")
case number < 50:
    fmt.Println("Medium")
default:
    fmt.Println("Large")
}
```

---

## 👉 Pointers

```go
// Create variable and pointer
x := 10
p := &x  // p points to x's address

// Print information
fmt.Println("Value:", x)        // Value: 10
fmt.Println("Address:", &x)     // Address: 0x... 
fmt.Println("Pointer:", p)      // Pointer: 0x...
fmt.Println("Dereferenced:", *p) // Dereferenced: 10

// Modify through pointer
*p = 20
fmt.Println("New value:", x)    // New value: 20

// Function with pointer parameter
func double(n *int) {
    *n = *n * 2
}

double(&x)  // Pass address of x
```

---

## 🏗️ Structs

```go
// Define struct
type Person struct {
    Name string
    Age  int
}

// Create instances
var p1 Person
p1.Name = "Alice"
p1.Age = 30

// Struct literal
p2 := Person{"Bob", 25}
p3 := Person{Name: "Charlie", Age: 35}

// Method with struct receiver argument
func (p Person) Greet() string {
    return fmt.Sprintf("Hi, I'm %s and I'm %d years old", p.Name, p.Age)
}

// Anonymous struct
temp := struct {
    ID   int
    Name string
}{
    ID:   1,
    Name: "Temp",
}
```

---

## 🚀 Quick Tips

### Common Patterns
```go
// Multiple assignment
a, b := 10, 20

// Swap variables  
a, b = b, a

// Ignore return value
_, err := someFunction()

// Type conversion
intValue := int(floatValue)
stringValue := string(runeValue)
```

### Best Practices
- Use `:=` for new variables
- Use `var` for zero values or explicit types
- Always handle the second return value from map lookups
- Use `range` for iteration when you don't need the index
- Name your return values for clarity in complex functions

### Memory
- **Arrays**: Fixed size, value type (copied when passed)
- **Slices**: Dynamic, reference type (header + underlying array)
- **Maps**: Reference type, always use `make()` or literal syntax
- **Pointers**: Hold memory addresses, use `&` to get address, `*` to dereference