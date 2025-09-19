package main

import (
	"fmt"
	"strconv"
)

// define person struct
type Person struct {
	firstName string
	lastName string 
	city string
	gender string
	age int
}

// type of methods/functions accessor/mutators 
//   called "Value Receiver, Pointer Receiver"

// Value Receiver
func (p Person) greet() string {
	return "Hello my name is " + p.firstName + " " + p.lastName +  " and I am "+ strconv.Itoa(p.age)
}

func (p *Person) hasBirthday() {
	p.age = p.age * 2
	
}




func main(){

	fmt.Println()

	person1 := Person {"Cristian", "M.", "Irvine", "M", 30}

	fmt.Println(person1)
	fmt.Println(person1.firstName)
	
	person1.age++

	fmt.Println(person1.greet())
	person1.hasBirthday()
	fmt.Println(person1.greet())
	
}