package user

import (
	"go-heroku-server/api/types"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserRole  = "user"
	AdminRole = "admin"
)

type User struct {
	ID        uint          `json:"-"`
	Username  string        `json:"username"`
	Password  string        `json:"-"`
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
	//err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	//if _, err = bcrypt.Cost([]byte(user.Password)); err != nil {
	//	log.Print(err)
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
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}
