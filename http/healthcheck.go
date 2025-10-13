package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Endpoint struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EndpointList struct {
	Endpoints []Endpoint `json:"endpoints"`
}

// Function 1: Check a single endpoint
func checkEndpoint(endpoint Endpoint, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(endpoint.URL)
	if err != nil {
		fmt.Printf("✗ %s [%s] - Error: %v\n", endpoint.Name, endpoint.URL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("✓ %s [%s] - OK\n", endpoint.Name, endpoint.URL)
	} else {
		fmt.Printf("✗ %s [%s] - Status: %d\n", endpoint.Name, endpoint.URL, resp.StatusCode)
	}
}

// Function 2: Get endpoints from URL and check all concurrently
func healthCheck(endpointsURL string) error {
	// Get endpoints from API
	resp, err := http.Get(endpointsURL)
	if err != nil {
		return fmt.Errorf("failed to get endpoints: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get endpoints, status: %d", resp.StatusCode)
	}

	// Parse JSON response
	var endpoints EndpointList
	if err := json.NewDecoder(resp.Body).Decode(&endpoints); err != nil {
		return fmt.Errorf("failed to decode endpoints: %v", err)
	}

	// Check all endpoints concurrently
	var wg sync.WaitGroup
	for _, endpoint := range endpoints.Endpoints {
		wg.Add(1)
		go checkEndpoint(endpoint, &wg)
	}
	wg.Wait()

	return nil
}

func HealthCheck() {
	// Get endpoints from an API
	endpointsURL := "https://api.example.com/endpoints"

	if err := healthCheck(endpointsURL); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
