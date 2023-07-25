package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserClient struct {
	token string
}

func NewUserClient() UserClient {
	return UserClient{
		token: "localhost:8080",
	}
}

func (r UserClient) loginUser() (string, error) {
	loginCredentials := map[string]string{
		"username": "admin",
		"password": "admin",
	}

	// Marshal the struct into a JSON byte slice
	jsonData, err := json.Marshal(loginCredentials)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return "", err
	}

	// Make the POST request
	url := "http://localhost:8080/user/login" // Replace this with the actual API endpoint URL
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body as a string
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response data:", err)
		return "", err
	}

	return string(responseData), nil
}
