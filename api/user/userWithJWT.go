package user

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go-heroku-server/api/src"
	"log"
	"net/http"
	"time"
)

var tokenEncodeString = []byte("Wow, much safe")
var tokenExpiration = time.Now().Local().Add(time.Second * 3600).Unix()

//Function verifies user if it exists and has valid login credentials
func LoginUser(w http.ResponseWriter, r *http.Request) {

	var credentials Credentials
	var serverToken Token

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	persistedUser, err := getUserByUsername(credentials.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("getUserByUsername error: %s\n", err)
		return
	}

	//if err = persistedUser.validatePassword(credentials.Password); err != nil {
	if !persistedUser.validatePassword(credentials.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		serverToken = createToken(persistedUser)
		w.Header().Set("Content-Type", "application/json")
		payload, _ := json.Marshal(serverToken)
		log.Printf("User %s logged successfully!\n", persistedUser.Username)
		_, _ = w.Write(payload)
	}
}

//TODO CustomClaims renaming to something like UserClaims
type CustomClaims struct {
	Id       uint
	Username string
	Role     string
	jwt.StandardClaims
}

//Function for creating token from verified user from LoginUser function
func createToken(verifiedUser User) (userToken Token) {
	customClaims := CustomClaims{
		Id:       verifiedUser.ID,
		Username: verifiedUser.Username,
		Role:     verifiedUser.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiration,
		},
	}

	var serverToken Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString(tokenEncodeString)
	serverToken.Token = tokenString

	return serverToken
}

// ParseToken Method extracts user CustomClaims from token
func ParseToken(signedToken string) (claims *CustomClaims, res src.IResponse) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return tokenEncodeString, nil
		},
	)

	if err != nil {
		res = src.NewErrorResponse(http.StatusUnauthorized, err)
		return
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		err = errors.New("could not parse claims")
		res = src.NewErrorResponse(http.StatusUnauthorized, err)
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT token has expired")
		res = src.NewErrorResponse(http.StatusUnauthorized, err)
		return
	}

	return
}
