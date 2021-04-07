package user

import "go-heroku-server/api/types"

type User struct {
	ID        uint          `json:"-"`
	Username  string        `json:"username"`
	Password  string        `json:"password"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Address   types.Address `json:"address"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}
