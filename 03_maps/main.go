package main

import (
	"fmt"
)

func main(){


	// define a map 
	emails := make(map[string]string)

	emails["Cristian"] = "xtian@test.local"
	emails["TestUser"] = "test@test.local"

	fmt.Println(emails)
	fmt.Println(len(emails))
	fmt.Println(emails["Cristian"])

	// delete from map
	delete(emails, "TestUser")
	fmt.Println(emails)

	// declare map and add kv
	emails2 := map[string]string {"Cristian":"xtian@test.local","TestUser":"test@test.local"}

	fmt.Println(emails2)
}