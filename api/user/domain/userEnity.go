package userDomains

import (
	"go-heroku-server/api/types"
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

func (u User) GetUsername() string {
	return u.Username
}

func (u User) GetRole() string {
	return u.Role
}

func (u *User) SetRole() {
	if u.Role == "" {
		u.Role = UserRole
	}
}
