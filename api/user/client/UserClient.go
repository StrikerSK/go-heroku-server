package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UserClient struct {
	BaserURL string
	Token    string
}

func NewUserClient() UserClient {
	userClient := UserClient{
		BaserURL: "http://localhost:8080",
	}
	userClient.fetchToken()
	return userClient
}

func (c *UserClient) fetchToken() {
	token, err := loginUser(c.BaserURL, "admin", "admin")
	if err != nil {
		log.Fatalf("error signing user: %v", err)
	}

	c.Token = token
}

func loginUser(baseURL, username, password string) (string, error) {
	loginCredentials := map[string]string{
		"username": username,
		"password": password,
	}

	// Marshal the struct into a JSON byte slice
	jsonData, err := json.Marshal(loginCredentials)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return "", err
	}

	// Make the POST request
	url := baseURL + "/user/login"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error making POST request:", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("user cannot be logged")
	}

	token := resp.Header.Get("Authorization")
	return token, nil
}
