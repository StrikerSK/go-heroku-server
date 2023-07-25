package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RetrievingUserToken(t *testing.T) {
	token, err := loginUser("admin", "admin")
	assert.Nil(t, err, "There should be no error during token retrieval")
	assert.NotEmpty(t, token, "Token should be returned")
}

func Test_InitializingUserClient(t *testing.T) {
	token := NewUserClient()
	assert.NotEmpty(t, token.token, "Token should be returned")
}
