package main

import (
	"fmt"
)

func PrintDataTypes() {
	// TODO: Declare variables of these types and assign values:
	// - string
	// - int
	// - float64
	// - bool
	// - rune (single character)

	// Print all variables with descriptive labels
	// Try both explicit declaration and short variable declaration (:=)
	fmt.Printf("Hello from %s", "exercise 1\n")
	hello := "hello world"
	fmt.Printf("This is a string: %s\n", hello)

	age := 22
	fmt.Printf("This is a int: %d\n", age)

	decimal := 22.22
	fmt.Printf("This is a float: %f\n", decimal)

	truthy := true
	fmt.Printf("This is a bool: %t\n", truthy)

	char := 'a'
	fmt.Printf("This is a run: %c\n", char)
}
