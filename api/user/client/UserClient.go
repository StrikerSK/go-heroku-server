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
	token, err := loginUser("admin", "admin")
	if err != nil {
		panic(err)
	}

	return UserClient{
		token: token,
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

	// Read the response body as a string
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response data:", err)
		return "", err
	}

	var mapStruct map[string]string
	err = json.Unmarshal(responseData, &mapStruct)
	if err != nil {
		fmt.Println("Error unmarshalling data:", err)
		return "", err
	}

	return mapStruct["token"], nil
}
