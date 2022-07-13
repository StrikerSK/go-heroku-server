package userDomains

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

func NewCredentials(username, password string) Credentials {
	return Credentials{
		Username: username,
		Password: password,
	}
}

func (c *Credentials) clearPassword() {
	c.Password = ""
}

func (c *Credentials) decryptPassword() {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("decriptPasswrod error: %s\n", err)
	}

	c.Password = string(encryptedPassword)
}

func (c *Credentials) validatePassword(password string) bool {
	//if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	//	log.Printf("validatePassword error: %s\n", err)
	//	return false
	//}
	return c.Password == password
}
