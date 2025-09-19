package main

import "fmt"

func main(){

	a := 5
	b := &a

	fmt.Println(a, b)
	fmt.Printf("%T\n", b)
	fmt.Printf("%T\n", a)

	// dereference pointer &
	fmt.Println(*&a)
	fmt.Println(*b) 

	// change value of a by way of changing value of pointer b
	*b = 10
	fmt.Println(a)
	fmt.Println(*&a)
}