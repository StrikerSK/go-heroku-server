package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserClient struct {
	Token string
}

func NewUserClient() UserClient {
	token, err := loginUser("admin", "admin")
	if err != nil {
		panic(err)
	}

	return UserClient{
		Token: token,
	}
}

func loginUser(username, password string) (string, error) {
	loginCredentials := map[string]string{
		"username": username,
		"password": password,
	}

	// Marshal the struct into a JSON byte slice
	jsonData, err := json.Marshal(loginCredentials)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return "", err
	}

	// Make the POST request
	url := "http://localhost:8080/user/login"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("user cannot be logged")
	}

	token := resp.Header.Get("Authorization")
	return token, nil
}
