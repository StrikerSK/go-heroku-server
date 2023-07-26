package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_RetrievingUserToken(t *testing.T) {
	token, err := loginUser("admin", "admin")
	assert.Nil(t, err, "There should be no error during token retrieval")
	assert.NotEmpty(t, token, "Token should be returned")
}

func Test_RetrievingWrongUserToken(t *testing.T) {
	_, err := loginUser("tester", "tester")
	assert.Error(t, err, "There should be error logging in")
	assert.EqualError(t, err, "user cannot be logged", "There should be error logging in")
}

func Test_InitializingUserClient(t *testing.T) {
	token := NewUserClient()
	assert.NotEmpty(t, token.Token, "Token should be returned")
}

func Test_RepeatingUserLogin(t *testing.T) {
	firstToken, err := loginUser("admin", "admin")
	assert.Nil(t, err, "There should be no error during token retrieval")
	assert.NotEmpty(t, firstToken, "Token should be returned")

	time.Sleep(1 * time.Second)

	secondToken, err := loginUser("admin", "admin")
	assert.Nil(t, err, "There should be no error during token retrieval")
	assert.NotEmpty(t, secondToken, "Token should be returned")
	assert.NotEqualf(t, secondToken, firstToken, "tokens should not match: expected %s to not match actual %s", firstToken, secondToken)
}
