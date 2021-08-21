package user

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	Id       uint
	Username string
	Role     string
	jwt.StandardClaims
}

type ClaimValues interface {
	GetID() uint
	GetUsername() string
	GetRole() string
}
