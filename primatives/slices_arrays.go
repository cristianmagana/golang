package main

import "fmt"

func Slices() {

	// TODO:
	// 1. Create an array of 5 integers
	example := [5]int{0, 1, 2, 3, 4}
	fmt.Println(example)

	// 2. Create a slice from the array
	slice := example[:4]
	fmt.Println(slice)

	// 3. Append elements to the slice
	slice = append(slice, 5)
	fmt.Println(slice)

	// 4. Create a slice of strings with your favorite colors
	colors := []string{"black", "red", "gold", "purple"}
	fmt.Println(colors)

	// 5. Use range to iterate and print all elements

	for i := range colors {
		fmt.Println(colors[i])
	}

	// 6. Show slice length and capacity
	fmt.Println(len(colors))
	fmt.Println(cap(colors))
}
