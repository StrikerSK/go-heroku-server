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

// Representation of the User model
type User struct {
	ID        uint          `json:"-"`
	Username  string        `json:"username" example:"JohnyDoe"`
	Password  string        `json:"password"` example:"Custom"
	FirstName string        `json:"firstName" example:"John"`
	LastName  string        `json:"lastName" example:"Doe"`
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

func (user *User) validatePassword(password string) bool {
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

// Login payload for user
type Credentials struct {
	Username string `json:"username" example:"JohnyDoe"`
	Password string `json:"password" example:"SecretPassword"`
}

type Token struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6MiwiVXNlcm5hbWUiOiJ0ZXN0ZXIiLCJSb2xlIjoidXNlciIsImV4cCI6MTYxODI2MzM5Mn0.sq5SMVq4Q2UqhKUXglDc8KJV0OlRq0N_GTbJmn0jzVY"`
}
