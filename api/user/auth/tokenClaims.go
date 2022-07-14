package userAuth

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	Username string
	Role     string
	jwt.StandardClaims
}
