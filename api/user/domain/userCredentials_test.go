package userDomains

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncryption(t *testing.T) {
	user := NewCredentials("admin", "admin")
	user.EncryptPassword()
	assert.NotEqual(t, "admin", user.Password, "expected password being encrypted")

	err := user.ValidatePassword("admin")
	assert.Nil(t, err, "expected no errors during password validation")

	err = user.ValidatePassword("tester")
	assert.NotNil(t, err, "expected error during password validation")
	assert.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error(), "expected error message does not match")
}

func TestUserPassword(t *testing.T) {
	user := User{
		UserCredentials: UserCredentials{
			Username: "admin",
			Password: "admin",
		},
		UserID:    1,
		FirstName: "Admin",
		LastName:  "Admin",
		Role:      "user",
		Address:   Address{},
	}

	user.EncryptPassword()
	assert.NotEqual(t, "admin", user.Password, "expected password being encrypted")

	err := user.ValidatePassword("admin")
	assert.Nil(t, err, "expected no errors during password validation")

	err = user.ValidatePassword("tester")
	assert.NotNil(t, err, "expected error during password validation")
	assert.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error(), "expected error message does not match")
}
