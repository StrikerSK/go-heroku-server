package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
)

var tokenEncodeString = []byte("Wow, much safe")

func DecodeToken(receivedToken string) (*jwt.Token, error) {
	return jwt.Parse(receivedToken, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return tokenEncodeString, nil
	})
}

//Function verifies user if it exists and has valid login credentials
func LoginUser(w http.ResponseWriter, r *http.Request) {

	var credentials Credentials
	var serverToken Token

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	persistedUser, err := getUserFromDB(credentials.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Print(err.Error())
		return
	}

	if credentials.Password == persistedUser.Password {

		serverToken = createToken(persistedUser)
		w.Header().Set("Content-Type", "application/json")
		payload, _ := json.Marshal(serverToken)
		log.Printf("User %s logged successfully!", persistedUser.Username)
		_, _ = w.Write(payload)

	} else {

		log.Printf("User %s not logged successfully!", persistedUser.Username)
		http.Error(w, "Invalid password", http.StatusUnauthorized)

	}
}

type CustomClaims struct {
	Id       uint
	Username string
	jwt.StandardClaims
}

//Function for creating token from verified user from LoginUser function
func createToken(verifiedUser User) (userToken Token) {
	customClaims := CustomClaims{
		Id:       verifiedUser.ID,
		Username: verifiedUser.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Second * 15).Unix(),
		},
	}

	var serverToken Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString(tokenEncodeString)
	serverToken.Token = tokenString

	return serverToken
}

//Method extracts user id from token
func GetIdFromToken(receivedToken string) (userId uint, err error) {

	token, err := DecodeToken(receivedToken)
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		log.Printf("Received id %d", claims.Id)
		return claims.Id, nil

	} else {

		return 0, err

	}
}

func ParseToken(signedToken string) (claims *CustomClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return tokenEncodeString, nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return
}
