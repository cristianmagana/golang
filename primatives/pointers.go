package main

import (
	"fmt"
)

func Pointers() {
	fmt.Println("Pointers")
	// TODO:
	// 1. Create an integer variable
	x := 10
	// 2. Create a pointer to that variable

	p := &x
	// 3. Print the value, address, and value through pointer

	fmt.Println("Value of x:", x)
	fmt.Println("Address of x:", &x)
	fmt.Println("Value stored in pointer p (address of x):", p)
	fmt.Println("Value pointed to by p (*p):", *p) // *p dereferences p to get the value at that address

	// 4. Modify the value through the pointer

	*p = 20 // Changes the value of x to 20
	fmt.Println("New value of x after modifying via pointer:", x)

	// 5. Create a function that modifies a value using pointers
	modifyPointer(p)
	fmt.Printf("Modified values: %d\n", *p)

}

func modifyPointer(number *int) {
	*number = *number * 2
}
