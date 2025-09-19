package main

import (
	"fmt"
)

func ControlFlow() {
	// TODO: Create a number guessing game that:
	// 1. Sets a target number (e.g., 42)
	target := 42
	// 2. Uses a for loop to simulate guesses from 1 to 50
	for guess := 1; guess <= 50; guess++ {
		fmt.Printf("Guess #%d: %d ", guess, guess)

		// 3. Use if-else to check if guess is too high, too low, or correct
		if guess < target {
			fmt.Print("- Too low! ")
		} else if guess > target {
			fmt.Print("- Too high! ")
		} else {
			fmt.Print("- CORRECT! ")
		}

		// 4. Use switch statement to categorize numbers
		switch {
		case guess%2 == 0:
			fmt.Print("(Even, ")
		default:
			fmt.Print("(Odd, ")
		}

		// Categorize by size
		switch {
		case guess <= 16:
			fmt.Print("Small)")
		case guess <= 33:
			fmt.Print("Medium)")
		default:
			fmt.Print("Large)")
		}

		fmt.Println() // New line

		// 5. Break the loop when target is found
		if guess == target {
			fmt.Printf("\n🎉 Found the target number %d in %d guesses!\n", target, guess)
			break
		}
	}
}
