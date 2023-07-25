package client

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		return "", nil
	}

	// Create the HTTP request with the file content in the body
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/user/login", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Println("Creating request:", err)
		return "", err
	}

	// Make the request to upload the file
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed with status code: %d", resp.StatusCode)
	}

	var response []byte
	_, _ = resp.Body.Read(response)

	return string(response), nil
}
