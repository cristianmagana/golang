package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// for large JSON response use "JSON-to-go" https://mholt.github.io/json-to-go/

type Token struct {
    Access_Token string `json:"access_token"`
    Token_type string `json:"token_type"`
    Expires_in string `json:"expires_in"`
 }

type Response struct {
	StatusCode int
 }

type Request struct {
	Response Response
	Token Token
	
}

func make_oauth_request() *Request {

	url := "https://somewebsite.com"
	payload := strings.NewReader("grant_type=client_credentials&scope=read write&client_id=XXXXXXX&client_secret=XXXXX")

	req, err := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "ASP.NET_SessionId=CookieSecret")
	
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

	var respToken Token 
	if err := json.Unmarshal([]byte(bodyBytes), &respToken); err != nil {
		fmt.Print(err.Error())
	}

	var response Response
	response.StatusCode = resp.StatusCode

	var request Request 
	request.Response = response
	request.Token = respToken

	return &request
}

