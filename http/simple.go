package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Address  Address `json:"address"`
	Phone    string  `json:"phone"`
	Website  string  `json:"website"`
	Company  Company `json:"company"`
}

type Address struct {
	Street  string `json:"street"`
	Suite   string `json:"suite"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
}

type Company struct {
	Name        string `json:"name"`
	CatchPhrase string `json:"catchPhrase"`
	BS          string `json:"bs"`
}

// Helper function to pretty print
func prettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("Error marshaling: %v", err)
		return
	}
	fmt.Println(string(b))
}

// GET request
func getUser(id int) (*User, error) {
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d", id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	return &user, nil
}

// POST request (create new user)
func createUser(user User) (*User, error) {
	url := "https://jsonplaceholder.typicode.com/users"

	// Marshal user to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %w", err)
	}

	// Make POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// Decode response
	var createdUser User
	if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	return &createdUser, nil
}

// PUT request (update existing user)
func updateUser(id int, user User) (*User, error) {
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d", id)

	// Marshal user to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %w", err)
	}

	// Create PUT request
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("PUT request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// Decode response
	var updatedUser User
	if err := json.NewDecoder(resp.Body).Decode(&updatedUser); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	return &updatedUser, nil
}

// DELETE request (bonus!)
func deleteUser(id int) error {
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d", id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("creating request failed: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("DELETE request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

func HttpSimple() {
	fmt.Println("=== GET User ===")
	user, err := getUser(1)
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint(user)

	fmt.Println("\n=== POST (Create) User ===")
	newUser := User{
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		Phone:    "123-456-7890",
		Website:  "johndoe.com",
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			Zipcode: "10001",
		},
		Company: Company{
			Name:        "Doe Inc",
			CatchPhrase: "Innovation at its best",
			BS:          "synergize efficient solutions",
		},
	}

	createdUser, err := createUser(newUser)
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint(createdUser)

	fmt.Println("\n=== PUT (Update) User ===")
	updatedUser := *user
	updatedUser.Name = "Updated Name"
	updatedUser.Email = "updated@example.com"

	result, err := updateUser(1, updatedUser)
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint(result)

	fmt.Println("\n=== DELETE User ===")
	if err := deleteUser(1); err != nil {
		log.Fatal(err)
	}
	fmt.Println("User deleted successfully!")
}
