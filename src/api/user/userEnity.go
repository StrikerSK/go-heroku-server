package user

import (
	"go-heroku-server/src/api/types"
	"golang.org/x/crypto/bcrypt"
	"log"
)

const (
	GuestRole = "guest"
	AdminRole = "admin"
)

type User struct {
	Credentials
	ID        uint          `json:"-"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Role      string        `json:"-" gorm:"default:guest"`
	Address   types.Address `json:"address"`
}

func (u User) GetID() uint {
	return u.ID
}

func (u User) GetUsername() string {
	return u.Username
}

func (u User) GetRole() string {
	return u.Role
}

func (u *User) decryptPassword() {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("decriptPasswrod error: %s\n", err)
	}

	u.Password = string(encryptedPassword)
}

func (u User) validatePassword(password string) bool {
	//if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	//	log.Printf("validatePassword error: %s\n", err)
	//	return false
	//}
	return u.Password == password
}

func (u *User) setRole() {
	if u.Role == "" {
		u.Role = GuestRole
	}
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

func (c *Credentials) clearPassword() {
	c.Password = ""
}

type Token struct {
	Token string `json:"token"`
}
