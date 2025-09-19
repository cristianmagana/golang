package main

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) Greet() string {
	return fmt.Sprintf("Hello my name is %s and I am %d years old\n", p.Name, p.Age)
}

func Structs() {
	fmt.Println("Structs")

	// TODO:
	// 1. Create a Person instance
	var doug Person
	doug.Name = "Doug"
	doug.Age = 44
	// 2. Create a Person using struct literal
	dougFunny := Person{"Doug", 44}
	fmt.Println("Person 2:", dougFunny.Greet())

	// 3. Create a method for the Person struct

	anonStruct := struct {
		Name string
		Age  int
	}{
		Name: "QuailMan",
		Age:  12,
	}
	fmt.Printf("Hello I am %s and I am %d years old", anonStruct.Name, anonStruct.Age)
	// 4. Use anonymous structs

}
