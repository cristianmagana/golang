package main

import (
	"fmt"
	"strings"
)

func StringManipulation(sentence string) {

	// TODO: Use the sentence variable to:
	// 1. Remove leading/trailing whitespace
	fmt.Println(strings.TrimSpace(sentence))
	// 2. Convert to uppercase
	fmt.Println(strings.ToUpper(sentence))
	// 3. Replace "World" with "Go"
	fmt.Println(strings.Replace(sentence, "World", "Go", -1))
	// 4. Split into words
	fmt.Println(strings.Split(sentence, ","))
	// 5. Check if it contains "Go"
	fmt.Println(strings.Contains(sentence, "Go"))
	// 6. Find the length of the original sentence
	fmt.Println(len(sentence))
	// Print results for each operation
}
