package userServices

import (
	"github.com/dgrijalva/jwt-go"
	"go-heroku-server/api/src/errors"
	userDomains "go-heroku-server/api/user/domain"
	"log"
	"time"
)

type TokenService struct {
	tokenEncoding   []byte
	tokenExpiration time.Duration
}

func NewTokenService(tokenEncoding string, tokenExpiry time.Duration) TokenService {
	return TokenService{
		tokenEncoding:   []byte(tokenEncoding),
		tokenExpiration: tokenExpiry,
	}
}

// Function for creating token from verified user from LoginUser function
func (s TokenService) CreateToken(user userDomains.User) (string, error) {
	claims := userDomains.UserClaims{
		UserID:   user.UserID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Second * s.tokenExpiration).Unix(),
			IssuedAt:  time.Now().Local().Unix(),
		},
	}

	// Sign and get the complete encoded token as a string using the secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.tokenEncoding)
	if err != nil {
		log.Printf("Create Token: %s\n", err.Error())
		return "", err
	}

	return signedToken, nil
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
		log.Printf("token parse error: %s\n", err.Error())
		err = errors.NewUnauthorizedError(err.Error())
		return
	}

	claims, ok := token.Claims.(*userDomains.UserClaims)
	if !ok {
		log.Printf("cannot resolve token claims")
		err = errors.NewUnauthorizedError("cannot resolve token claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		log.Printf("JWT token has expired")
		err = errors.NewUnauthorizedError("JWT token has expired")
		return
	}

	return
}
