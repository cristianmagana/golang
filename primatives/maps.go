package main

import (
	"fmt"
)

func Maps() {
	// TODO: Create a program that:
	// 1. Creates a map of student names to their grades
	grades := make(map[string]string)
	fmt.Println(grades)

	// 2. Adds several students
	grades["Emilio"] = "C"
	grades["Zoe"] = "A"
	grades["Claudia"] = "F"
	fmt.Println(grades)
	// 3. Updates a student's grade
	grades["Claudia"] = "A"
	fmt.Println(grades)

	// 4. Checks if a student exists

	value, ok := grades["Zoe"]
	if ok {
		fmt.Println(value)
	} else {
		fmt.Println("DNE")
	}

	// 5. Deletes a student
	delete(grades, "Zoe")

	// 6. Iterates through all students and prints them
	for key, value := range grades {
		fmt.Printf("key: %s, value: %s\n", key, value)
	}
}
