package userDomains

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

func NewCredentials(username, password string) UserCredentials {
	return UserCredentials{
		Username: username,
		Password: password,
	}
}

func (c *UserCredentials) ClearPassword() {
	c.Password = ""
}

func (c *UserCredentials) DecryptPassword() {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("decriptPasswrod error: %s\n", err)
	}

	c.Password = string(encryptedPassword)
}

func (c *UserCredentials) ValidatePassword(password string) bool {
	//if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	//	log.Printf("validatePassword error: %s\n", err)
	//	return false
	//}
	return c.Password == password
}
