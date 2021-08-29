package main

import "fmt"

func main(){

	ids := []int{11,22,33,44,55,66,77,88,99.00}

	// loop through ids
	for i, id := range ids {
		fmt.Printf("%d - ID: %d\n", i, id)
	}

	// not using index, using blank identifier
	for _, id := range ids {
		fmt.Printf("ID: %d\n", id)
	}

	// range with map 
	emails2 := map[string]string {"Cristian":"xtian@test.local","TestUser":"test@test.local"}

	for k, v := range emails2 {
		fmt.Printf("%s: %s\n", k, v)
	}

	for k := range emails2 {
		fmt.Println("Name: " + k)
	}
}