package user

import (
	uuid "github.com/satori/go.uuid"
	"go-heroku-server/api/src"
	"go-heroku-server/config"
	"net/http"
	"time"
)

const tokenName = "session_token"

func login(credentials Credentials) (*http.Cookie, *src.RequestError) {
	var requestError src.RequestError

	persistedUser, err := getUserByUsername(credentials.Username)
	if err != nil {
		requestError = src.NewErrorResponse(http.StatusBadRequest, err)
		return nil, &requestError
	}

	// If a password exists for the given user
	// AND, if it is the same as the password we received, then we can move ahead
	// if NOT, then we return an "Unauthorized" status
	//if err = persistedUser.validatePassword(credentials.Password); err != nil {
	if !persistedUser.validatePassword(credentials.Password) {
		requestError = src.NewErrorResponse(http.StatusUnauthorized, err)
		return nil, &requestError
	}

	// Create a new random session token
	sessionToken := uuid.NewV4().String()
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	_, err = config.Cache.Do("SETEX", sessionToken, "120", persistedUser.ID)
	if err != nil {
		// If there is an error in setting the cache, return an internal server error
		requestError = src.NewErrorResponse(http.StatusInternalServerError, err)
		return nil, &requestError
	}

	cookie := &http.Cookie{
		Name:    tokenName,
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	}

	return cookie, nil
}
