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
	ID        uint          `json:"-"`
	Username  string        `json:"username"`
	Password  string        `json:"password"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Role      string        `json:"-" gorm:"default:user"`
	Address   types.Address `json:"address"`
}

func (user *User) decryptPassword() {
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(encryptedPassword)
}

func (user *User) validatePassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Print(err)
		return false
	}
	return true
}

func (user *User) setRole() {
	if user.Role == "" {
		user.Role = UserRole
	}
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}
