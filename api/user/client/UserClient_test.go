package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetrievingUserToken(t *testing.T) {
	login := NewUserClient()
	token, err := login.loginUser()
	assert.Nil(t, err, "There should be no error during token retrieval")
	assert.NotNil(t, token, "Token should be returned")
}
