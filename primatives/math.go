package main

import (
	"fmt"
	"math"
)

func Math() {

	// TODO: Create a calculator that:
	// 1. Declares two integers (a=10, b=3)
	a, b := 5, 10
	fmt.Printf("The two variables are: 1: %d, 2: %d\n", a, b)
	// 2. Performs and prints: addition, subtraction, multiplication, division, modulo
	fmt.Println(a + b)
	fmt.Println(b - a)
	fmt.Println(a * b)
	fmt.Println(b % a)
	// 3. Demonstrates integer vs float division
	fmt.Println(b / a)
	fmt.Println(float64(b) / float64(a))
	// 4. Uses math package functions: Sqrt, Pow, Ceil, Floor
	fmt.Println(math.Sqrt(64))
	fmt.Println(math.Pow(float64(a), 2))
	fmt.Println(math.Ceil(10.1))
	fmt.Println(math.Floor(10.1))
	// Bonus: Show type conversion between int and float64

}
