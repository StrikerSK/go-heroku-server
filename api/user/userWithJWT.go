package user

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"go-heroku-server/api/src/responses"
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
		log.Printf("Logging error: %s\n", err)
		responses.NewEmptyResponse(http.StatusBadRequest).WriteResponse(w)
		return
	}

	persistedUser, err := getUserByUsername(credentials.Username)
	if err != nil {
		log.Printf("Logging error: %s\n", err)
		responses.NewEmptyResponse(http.StatusUnauthorized).WriteResponse(w)
		return
	}

	//if err = persistedUser.validatePassword(credentials.Password); err != nil {
	if !persistedUser.validatePassword(credentials.Password) {
		responses.NewEmptyResponse(http.StatusUnauthorized)
		return
	} else {
		serverToken = createToken(persistedUser)
		log.Printf("User [%d] login: success\n", persistedUser.ID)
		responses.NewResponse(serverToken).WriteResponse(w)
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
func ParseToken(signedToken string) (claims *CustomClaims, res responses.IResponse) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return tokenEncodeString, nil
		},
	)

	if err != nil {
		log.Printf("Token parse: %s\n", err.Error())
		res = responses.NewEmptyResponse(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		log.Printf("Token parse: %s\n", err.Error())
		res = responses.NewEmptyResponse(http.StatusUnauthorized)
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		log.Println("Token parse: JWT token has expired")
		res = responses.NewEmptyResponse(http.StatusUnauthorized)
		return
	}

	return
}
