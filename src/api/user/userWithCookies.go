package user

import (
	uuid "github.com/satori/go.uuid"
	"go-heroku-server/config"
	"go-heroku-server/src/responses"
	"log"
	"net/http"
	"time"
)

const tokenName = "session_token"

func loginWithCookies(credentials Credentials) (http.Cookie, responses.IResponse) {
	var requestError responses.IResponse

	persistedUser, err := validateUser(credentials)
	if err != nil {
		log.Printf("Login user(Cookie): %v\n", err)
		requestError = responses.CreateResponse(http.StatusBadRequest, nil)
		return http.Cookie{}, requestError
	}

	// Create a new random session token
	sessionToken := uuid.NewV4().String()
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 3600 seconds
	if _, err = config.GetCacheInstance().Do("SETEX", sessionToken, "3600", persistedUser.ID); err != nil {
		// If there is an error in setting the cache, return an internal server error
		log.Printf("Login user(Cookie): %v\n", err)
		requestError = responses.CreateResponse(http.StatusInternalServerError, nil)
		return http.Cookie{}, requestError
	}

	cookie := http.Cookie{
		Name:    tokenName,
		Value:   sessionToken,
		Expires: time.Now().Add(3600 * time.Second),
	}

	return cookie, nil
}
