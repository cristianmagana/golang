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

	// array
	arrayEx := [5]string{"array","of","fixed","size","0-based"}
	fmt.Println(arrayEx)

	// slice i.e.; vector
	sliceEx := []string{"c++","like","vector"}
	fmt.Println(sliceEx)

	fmt.Println(name, age, isMale)

}