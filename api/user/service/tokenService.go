package userServices

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	userDomains "go-heroku-server/api/user/domain"
	"log"
	"time"
)

type TokenService struct {
	tokenEncoding   []byte
	tokenExpiration int64
}

func NewTokenService(tokenEncoding string, tokenExpiry time.Duration) TokenService {
	return TokenService{
		tokenEncoding:   []byte(tokenEncoding),
		tokenExpiration: time.Now().Local().Add(time.Second * tokenExpiry).Unix(),
	}
}

//Function for creating token from verified user from LoginUser function
func (s TokenService) CreateToken(user userDomains.User) (token string, err error) {
	claims := userDomains.UserClaims{
		UserID:   user.UserID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: s.tokenExpiration,
		},
	}

	// Sign and get the complete encoded token as a string using the secret
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.tokenEncoding)
	if err != nil {
		log.Printf("Create Token: %s\n", err.Error())
		return
	}

	return
}

// ParseToken Method extracts user CustomClaims from token
func (s TokenService) ParseToken(signedToken string) (claims *userDomains.UserClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&userDomains.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return s.tokenEncoding, nil
		},
	)

	if err != nil {
		log.Printf("Token parse: %s\n", err.Error())
		return
	}

	claims, ok := token.Claims.(*userDomains.UserClaims)
	if !ok {
		err = errors.New("cannot resolve token claims")
		log.Printf("Token parse: %s\n", err.Error())
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT token has expired")
		log.Printf("Token parse: %s\n", err.Error())
		return
	}

	return
}
