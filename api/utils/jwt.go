package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go-heroku-server/api/user"
	"strconv"
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

//Function for creating token from verified user from LoginUser function
func CreateToken(verifiedUser user.User) (userToken user.Token) {

	var serverToken user.Token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       fmt.Sprint(verifiedUser.ID),
		"username": verifiedUser.Username,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString(tokenEncodeString)
	serverToken.Token = tokenString

	return serverToken

}

//func VerifyToken(w http.ResponseWriter, r *http.Request) {
//
//	var loginCredentials UserCredentials;
//	json.NewDecoder(r.Body).Decode(&loginCredentials)
//
//	receivedToken := (r.Header.Get("Authorization"))
//	fmt.Println("Received token: " + receivedToken)
//
//	token, err := DecodeToken(receivedToken)
//
//	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//
//		userId := claims["id"].(string)
//		newUserId, _ := strconv.ParseUint(userId, 10, 32)
//
//		//userAddress := GetUserAddress(uint(newUserId))
//
//		json.NewEncoder(w).Encode(userAddress)
//
//	} else {
//
//		fmt.Println(err)
//	}
//}

//Method extracts user id from token
func GetIdFromToken(receivedToken string) (userId uint, err error) {

	token, err := DecodeToken(receivedToken)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		userId := claims["id"].(string)
		newUserId, _ := strconv.ParseUint(userId, 10, 32)

		return uint(newUserId), nil

	} else {

		return 0, err

	}
}
