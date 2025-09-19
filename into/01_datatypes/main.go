package main

import (
	"fmt"
)

func main() {

	// string
	// bool
	// int int8 int16 int32 int64
	// uint ""
	// byte / uint8
    // rune / int32
	// float32 float64
	// complex64 complex128

	// Datatypes are infered upon declaration. 
	name := "Cristian"
	age  := 20
	isMale := 1
	fmt.Println(name, age, isMale)

	// array
	arrayEx := [5]string{"array","of","fixed","size","0-based"}
	fmt.Println(arrayEx)

	// slice i.e.; vector
	sliceEx := []string{"c++","like","vector"}
	fmt.Println(sliceEx)

	// conditionals
	color := "red"

	if color == "red" {
		fmt.Println("is red")
	} else if color == "blue" {
		fmt.Println("is blue")
	} else {
		fmt.Println("not blue nor red")
	}

	// switch 
	switch color {
	case "red":
		fmt.Println("is red")
	case "blue":
		fmt.Println("is blue")
	default:
		fmt.Println("not blue nor red")

	}
}