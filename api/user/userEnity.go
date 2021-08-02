package user

import (
	"go-heroku-server/api/types"
	"golang.org/x/crypto/bcrypt"
	"log"
)

const (
	UserRole  = "user"
	AdminRole = "admin"
)

type User struct {
	Credentials
	ID        uint          `json:"-"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Role      string        `json:"-" gorm:"default:user"`
	Address   types.Address `json:"address"`
}

func (user *User) decryptPassword() {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("decriptPasswrod error: %s\n", err)
	}

	user.Password = string(encryptedPassword)
}

func (user User) validatePassword(password string) bool {
	//if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	//	log.Printf("validatePassword error: %s\n", err)
	//	return false
	//}
	return user.Password == password
}

func (user *User) setRole() {
	if user.Role == "" {
		user.Role = UserRole
	}
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type Token struct {
	Token string `json:"token"`
}
