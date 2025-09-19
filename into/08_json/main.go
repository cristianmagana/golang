package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type JsonStruct struct {

	Fname string `json:"fname"`
	Lname string `json:"lname"`
	DoB int `json:"dob"`
}


func main () {


	// Unmarshal into struct
	req, err := http.NewRequest("POST", "http://example.com", nil)

	if err != nil {
		fmt.Print(err.Error())
	}

	client := &http.Client{}
	
	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var respJson JsonStruct
	if err := json.Unmarshal([]byte(bodyBytes), &respJson); err != nil {
		fmt.Print(err.Error())
	}

	fmt.Println(respJson)


	fmt.Println()
}