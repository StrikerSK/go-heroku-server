package client

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type UserClient struct {
	BaserURL string
	Token    string
}

func NewUserClient() UserClient {
	baseURL := "http://localhost:8080"

	userClient := UserClient{
		BaserURL: baseURL,
		Token:    loginUser(baseURL, "admin", "admin"),
	}
	//userClient.fetchToken()
	return userClient
}

//func (c *UserClient) fetchToken() {
//	token, err := loginUser(c.BaserURL, "admin", "admin")
//	if err != nil {
//		log.Fatalf("error signing user: %v", err)
//	}
//
//	c.Token = token
//}

func loginUser(baseURL, username, password string) string {
	loginCredentials := map[string]string{
		"username": username,
		"password": password,
	}

	// Marshal the struct into a JSON byte slice
	jsonData, err := json.Marshal(loginCredentials)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
	}

	// Make the POST request
	url := baseURL + "/user/login"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error making POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		log.Println("user cannot be logged")
		panic("user cannot be logged")
	}

	token := resp.Header.Get("Authorization")
	return token
}
