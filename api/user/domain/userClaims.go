package userDomains

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	UserID   uint
	Username string
	Role     string
	jwt.StandardClaims
}
