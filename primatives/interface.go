package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ============================================================================
// WHAT IS interface{} ?
// ============================================================================
// interface{} is the EMPTY INTERFACE - it has zero methods
// Since every type implements zero methods, EVERY type satisfies interface{}
// Think of it as: "I'll accept ANY type"

// Go 1.18+ introduced 'any' as an alias for interface{}
// any == interface{}  (they're identical)

func demonstrateEmptyInterface() {
	fmt.Println("=== Empty Interface Basics ===\n")

	var x interface{} // Can hold ANY type

	x = 42
	fmt.Printf("x = %v (type: %T)\n", x, x) // x = 42 (type: int)

	x = "hello"
	fmt.Printf("x = %v (type: %T)\n", x, x) // x = hello (type: string)

	x = []int{1, 2, 3}
	fmt.Printf("x = %v (type: %T)\n", x, x) // x = [1 2 3] (type: []int)

	x = map[string]int{"age": 30}
	fmt.Printf("x = %v (type: %T)\n\n", x, x) // x = map[age:30] (type: map[string]int)
}

// ============================================================================
// COMPARISON: Go interface{} vs TypeScript any
// ============================================================================

func compareWithTypeScript() {
	fmt.Println("=== Go interface{} vs TypeScript any ===\n")

	fmt.Println(`
SIMILARITIES:
✓ Both accept any type
✓ Both require type checking/assertion to use
✓ Both lose compile-time type safety
✓ Both used for dynamic/unknown data

DIFFERENCES:

Go interface{}:
  - Runtime type information preserved
  - Type assertions can fail at runtime (panic or return false)
  - No implicit operations (must type assert first)
  - Value is stored with type metadata (slight overhead)
  
TypeScript any:
  - Compiles to JavaScript (no runtime types)
  - Type checking only at compile time
  - Can perform any operation (unsafe!)
  - No runtime overhead (it's just JavaScript)

Go Example:
  var x interface{} = "hello"
  s := x.(string)           // Type assertion - can panic!
  s, ok := x.(string)       // Safe type assertion
  
TypeScript Example:
  let x: any = "hello"
  let s: string = x         // No runtime check!
  s.toUpperCase()           // Works, but could crash if x isn't string
`)
}

// ============================================================================
// TYPE ASSERTIONS - How to Use interface{} Values
// ============================================================================

func typeAssertions() {
	fmt.Println("=== Type Assertions ===\n")

	var val interface{} = 42

	// Method 1: Direct assertion (panics if wrong type)
	fmt.Println("Method 1: Direct assertion")
	num := val.(int)
	fmt.Printf("Value: %d\n", num)

	// This would panic!
	// str := val.(string) // panic: interface conversion: interface {} is int, not string

	// Method 2: Safe assertion with ok pattern (recommended)
	fmt.Println("\nMethod 2: Safe assertion (comma-ok)")
	if str, ok := val.(string); ok {
		fmt.Printf("It's a string: %s\n", str)
	} else {
		fmt.Println("Not a string!")
	}

	if num, ok := val.(int); ok {
		fmt.Printf("It's an int: %d\n", num)
	}

	// Method 3: Type switch (best for multiple types)
	fmt.Println("\nMethod 3: Type switch")
	switch v := val.(type) {
	case int:
		fmt.Printf("Integer: %d\n", v)
	case string:
		fmt.Printf("String: %s\n", v)
	case []interface{}:
		fmt.Printf("Slice: %v\n", v)
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
	fmt.Println()
}

// ============================================================================
// map[string]interface{} - Dynamic JSON-like Data
// ============================================================================

func mapStringInterface() {
	fmt.Println("=== map[string]interface{} - The JSON Pattern ===\n")

	// This is Go's equivalent to TypeScript's { [key: string]: any }
	// Commonly used for:
	// 1. Parsing JSON with unknown structure
	// 2. Dynamic configuration
	// 3. API responses with variable fields

	data := map[string]interface{}{
		"name":    "Alice",
		"age":     30,
		"active":  true,
		"scores":  []int{95, 87, 92},
		"address": map[string]string{"city": "NYC"},
	}

	fmt.Println("Data structure:")
	for key, val := range data {
		fmt.Printf("  %s: %v (type: %T)\n", key, val, val)
	}

	// Accessing values requires type assertion
	fmt.Println("\nAccessing values:")

	// Safe access with type assertion
	if name, ok := data["name"].(string); ok {
		fmt.Printf("Name: %s\n", name)
	}

	if age, ok := data["age"].(int); ok {
		fmt.Printf("Age: %d\n", age)
	}

	// Accessing nested data
	if addr, ok := data["address"].(map[string]string); ok {
		fmt.Printf("City: %s\n", addr["city"])
	}

	fmt.Println()
}

// ============================================================================
// JSON Unmarshaling - Primary Use Case
// ============================================================================

func jsonExample() {
	fmt.Println("=== JSON Unmarshaling Example ===\n")

	// Unknown JSON structure
	jsonData := `{
		"user": "bob",
		"count": 42,
		"metadata": {
			"created": "2025-01-01",
			"tags": ["go", "tutorial"]
		}
	}`

	// Option 1: Unmarshal to map[string]interface{}
	var result map[string]interface{}
	json.Unmarshal([]byte(jsonData), &result)

	fmt.Println("Parsed JSON:")
	for key, val := range result {
		fmt.Printf("  %s: %v (type: %T)\n", key, val, val)
	}

	// Accessing nested JSON
	if metadata, ok := result["metadata"].(map[string]interface{}); ok {
		if tags, ok := metadata["tags"].([]interface{}); ok {
			fmt.Println("\nTags:")
			for i, tag := range tags {
				fmt.Printf("  [%d]: %v\n", i, tag)
			}
		}
	}

	fmt.Println()
}

// ============================================================================
// BETTER ALTERNATIVES to interface{}
// ============================================================================

func betterAlternatives() {
	fmt.Println("=== Better Alternatives to interface{} ===\n")

	fmt.Println(`
1. STRUCT WITH KNOWN FIELDS (Best for known structure):

   type User struct {
       Name   string
       Age    int
       Active bool
   }
   
   ✓ Type-safe
   ✓ No type assertions needed
   ✓ Better performance
   ✓ IDE autocomplete

2. GENERICS (Go 1.18+) (Best for reusable code):

   func Max[T comparable](a, b T) T {
       if a > b { return a }
       return b
   }
   
   ✓ Type-safe at compile time
   ✓ No interface{} overhead
   ✓ Reusable across types

3. SPECIFIC INTERFACES (Best for behavior):

   type Reader interface {
       Read(p []byte) (n int, err error)
   }
   
   ✓ Documents expectations
   ✓ Type-safe
   ✓ Enables polymorphism

4. TYPE ALIASES (Best for clarity):

   type UserID string
   type Metadata map[string]string
   
   ✓ Self-documenting
   ✓ Type-safe
   ✓ Can add methods

WHEN TO USE interface{}/any:
  ✗ Avoid in function parameters (use generics instead)
  ✗ Avoid in struct fields (use concrete types)
  ✓ OK for JSON unmarshaling unknown structures
  ✓ OK for reflection-based libraries
  ✓ OK for fmt.Printf and similar variadic functions
`)
}

// ============================================================================
// PRACTICAL EXAMPLE: Configuration Parser
// ============================================================================

type Config struct {
	Settings map[string]interface{}
}

func (c *Config) GetString(key string) (string, error) {
	val, exists := c.Settings[key]
	if !exists {
		return "", fmt.Errorf("key not found: %s", key)
	}

	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("key %s is not a string (got %T)", key, val)
	}

	return str, nil
}

func (c *Config) GetInt(key string) (int, error) {
	val, exists := c.Settings[key]
	if !exists {
		return 0, fmt.Errorf("key not found: %s", key)
	}

	// JSON numbers unmarshal as float64!
	switch v := val.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("key %s is not a number (got %T)", key, val)
	}
}

func configExample() {
	fmt.Println("=== Practical Config Parser ===\n")

	jsonConfig := `{
		"app_name": "MyApp",
		"port": 8080,
		"debug": true
	}`

	var settings map[string]interface{}
	json.Unmarshal([]byte(jsonConfig), &settings)

	config := Config{Settings: settings}

	// Type-safe getters hide the interface{} complexity
	if appName, err := config.GetString("app_name"); err == nil {
		fmt.Printf("App Name: %s\n", appName)
	}

	if port, err := config.GetInt("port"); err == nil {
		fmt.Printf("Port: %d\n", port)
	}

	// Error handling for wrong types
	if _, err := config.GetString("port"); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println()
}

// ============================================================================
// REFLECTION - Advanced interface{} Usage
// ============================================================================

func reflectionExample() {
	fmt.Println("=== Reflection with interface{} ===\n")

	var val interface{} = []int{1, 2, 3}

	// Get runtime type information
	t := reflect.TypeOf(val)
	v := reflect.ValueOf(val)

	fmt.Printf("Type: %v\n", t)
	fmt.Printf("Kind: %v\n", t.Kind())
	fmt.Printf("Value: %v\n", v)

	// Check if it's a slice
	if t.Kind() == reflect.Slice {
		fmt.Printf("Slice length: %d\n", v.Len())
		fmt.Printf("Element type: %v\n", t.Elem())

		// Iterate over slice elements
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			fmt.Printf("  [%d]: %v\n", i, elem.Interface())
		}
	}

	fmt.Println()
}

// ============================================================================
// PERFORMANCE CONSIDERATIONS
// ============================================================================

func performanceNotes() {
	fmt.Println("=== Performance Considerations ===\n")

	fmt.Println(`
OVERHEAD OF interface{}:

1. Memory:
   - interface{} stores value + type information (2 words: 16 bytes on 64-bit)
   - Concrete type: just the value
   - Small values may be heap allocated when stored in interface{}

2. Performance:
   - Type assertions have small runtime cost
   - Type switches have small runtime cost
   - Prevents compiler optimizations
   - May cause heap allocations

BENCHMARK RESULTS (approximate):

  Direct int access:        1 ns/op
  interface{} type assert:  5 ns/op (5x slower)
  Type switch:              3 ns/op (3x slower)
  Reflection:              50 ns/op (50x slower!)

RULE OF THUMB:
  - Use concrete types for hot paths
  - interface{} is fine for I/O, JSON, config (already slow operations)
  - Avoid interface{} in tight loops
  - Use generics instead of interface{} when possible
`)
}

// ============================================================================
// REAL-WORLD PATTERNS
// ============================================================================

func realWorldPatterns() {
	fmt.Println("=== Real-World Patterns ===\n")

	fmt.Println(`
PATTERN 1: Generic Response Wrapper

type APIResponse struct {
    Success bool
    Data    interface{}  // Different endpoints return different types
    Error   string
}

// Usage:
return APIResponse{Success: true, Data: user}
return APIResponse{Success: true, Data: []Product{...}}

PATTERN 2: Event System

type Event struct {
    Type    string
    Payload interface{}  // Event-specific data
}

PATTERN 3: Caching

type Cache struct {
    data map[string]interface{}  // Cache any type
    mu   sync.RWMutex
}

PATTERN 4: Middleware Context

type Context struct {
    values map[string]interface{}  // Store arbitrary data
}

ANTI-PATTERNS (Avoid!):

❌ func Process(data interface{}) interface{}
   Use generics or specific types instead

❌ type User struct {
       Data interface{}  // What is this?
   }
   Use concrete fields instead

❌ Overusing interface{} everywhere
   Use it sparingly, prefer type safety
`)
}

// ============================================================================
// MODERN GO: 'any' keyword (Go 1.18+)
// ============================================================================

func modernAny() {
	fmt.Println("=== Modern Go: 'any' keyword ===\n")

	// 'any' is an alias for interface{} (they're identical)
	var x any = "hello"
	var y interface{} = "world"

	fmt.Printf("x type: %T, y type: %T\n", x, y)
	fmt.Println("Both are exactly the same thing!\n")

	// Modern Go code prefers 'any' for readability
	data := map[string]any{
		"name": "Bob",
		"age":  30,
	}

	fmt.Printf("Modern style with 'any': %v\n", data)

	// But interface{} still works and is common in older code
	oldStyle := map[string]interface{}{
		"name": "Alice",
		"age":  25,
	}

	fmt.Printf("Old style with interface{}: %v\n\n", oldStyle)
}

// ============================================================================
// MAIN
// ============================================================================

func Interface() {
	demonstrateEmptyInterface()
	compareWithTypeScript()
	typeAssertions()
	mapStringInterface()
	jsonExample()
	betterAlternatives()
	configExample()
	reflectionExample()
	performanceNotes()
	realWorldPatterns()
	modernAny()

	//fmt.Println("=" + "=".repeat(79))
	fmt.Println("KEY TAKEAWAYS:")
	//fmt.Println("=" + "=".repeat(79))
	fmt.Println(`
1. interface{} (or 'any') accepts ANY type - similar to TypeScript's any
2. You MUST type assert to use the value (Go is stricter than TS)
3. map[string]interface{} is the Go way to handle dynamic JSON
4. Use concrete types when possible - interface{} loses type safety
5. Modern Go: prefer 'any' over interface{} for readability
6. Common uses: JSON parsing, reflection, generic containers
7. Avoid in hot paths - has performance overhead
`)
}
