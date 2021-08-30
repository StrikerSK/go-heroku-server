package user

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

var tokenEncodeString = []byte("Wow, much safe")
var tokenExpiration = time.Now().Local().Add(time.Second * 3600).Unix()

//Function for creating token from verified user from LoginUser function
func CreateToken(user ClaimValues) (token string, err error) {
	claims := UserClaims{
		Id:       user.GetID(),
		Username: user.GetUsername(),
		Role:     user.GetRole(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiration,
		},
	}

	// Sign and get the complete encoded token as a string using the secret
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(tokenEncodeString)
	return
}

// ParseToken Method extracts user CustomClaims from token
func ParseToken(signedToken string) (claims *UserClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return tokenEncodeString, nil
		},
	)

	if err != nil {
		log.Printf("Token parse: %v\n", err)
		return
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		err = errors.New("cannot resolve token claims")
		log.Printf("Token parse: %v\n", err)
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT token has expired")
		log.Printf("Token parse: %v\n", err)
		return
	}

	return
}
