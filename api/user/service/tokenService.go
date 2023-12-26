package userServices

import (
	"github.com/dgrijalva/jwt-go"
	"go-heroku-server/api/src/errors"
	userDomains "go-heroku-server/api/user/domain"
	"go-heroku-server/config"
	"log"
	"time"
)

type TokenService struct {
	tokenEncoding   []byte
	tokenExpiration time.Duration
}

func NewTokenService(configuration config.Authorization) TokenService {
	return TokenService{
		tokenEncoding:   []byte(configuration.TokenEncoding),
		tokenExpiration: time.Duration(configuration.TokenExpiration),
	}
}

// Function for creating token from verified user from LoginUser function
func (s TokenService) CreateToken(user userDomains.User) (string, error) {
	localTime := time.Now().Local()
	claims := userDomains.UserClaims{
		UserID:   user.UserID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: localTime.Add(time.Second * s.tokenExpiration).Unix(),
			IssuedAt:  localTime.Unix(),
		},
	}

	// Sign and get the complete encoded token as a string using the secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.tokenEncoding)
	if err != nil {
		log.Printf("cannot create token: %s\n", err.Error())
		return "", err
	}

	return signedToken, nil
}

// ParseToken Method extracts user CustomClaims from token
func (s TokenService) ParseToken(signedToken string) (*userDomains.UserClaims, error) {
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
		return nil, err
	}

	claims, ok := token.Claims.(*userDomains.UserClaims)
	if !ok {
		log.Printf("cannot resolve token claims")
		err = errors.NewUnauthorizedError("cannot resolve token claims")
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		log.Printf("JWT token has expired")
		err = errors.NewUnauthorizedError("JWT token has expired")
		return nil, err
	}

	return claims, nil
}
