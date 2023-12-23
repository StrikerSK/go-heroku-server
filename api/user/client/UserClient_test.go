package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var baseURL = "http://localhost:8080"

func Test_RetrievingUserToken(t *testing.T) {
	token := loginUser(baseURL, "admin", "admin")
	assert.NotEmpty(t, token, "Token should be returned")
}

func Test_RetrievingWrongUserToken(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but no panic occurred")
		} else if r != "user cannot be logged" {
			t.Errorf("Expected panic message 'user cannot be logged', but got: %v", r)
		}
	}()
	loginUser(baseURL, "admin", "wrong")
}

func Test_InitializingUserClient(t *testing.T) {
	token := NewUserClient(baseURL, "admin", "admin")
	assert.NotEmpty(t, token.Token, "Token should be returned")
	assert.NotEmpty(t, token.BaserURL, "URL should not be empty")
}

func Test_InitializingUserClient_WrongCredentials(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic, but no panic occurred")
		} else if r != "user cannot be logged" {
			t.Errorf("Expected panic message 'user cannot be logged', but got: %v", r)
		}
	}()
	NewUserClient(baseURL, "admin", "wrong")
}

func Test_RepeatingUserLogin(t *testing.T) {
	firstToken := loginUser(baseURL, "admin", "admin")
	assert.NotEmpty(t, firstToken, "Token should be returned")

	time.Sleep(1 * time.Second)
	secondToken := loginUser(baseURL, "admin", "admin")
	assert.NotEmpty(t, secondToken, "Token should be returned")
	assert.NotEqualf(t, secondToken, firstToken, "tokens should not match: expected %s to not match actual %s", firstToken, secondToken)
}
