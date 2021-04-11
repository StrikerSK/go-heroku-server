package user

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"go-heroku-server/config"
	"log"
	"net/http"
	"time"
)

const tokenName = "session_token"

func login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err.Error())
		return
	}

	// Get the expected password from our in memory map
	persistedUser, err := getUserByUsername(credentials.Username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err.Error())
		return
	}

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	//if err = persistedUser.validatePassword(credentials.Password); err != nil {
	if !persistedUser.validatePassword(credentials.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a new random session token
	sessionToken := uuid.NewV4().String()
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	_, err = config.Cache.Do("SETEX", sessionToken, "120", persistedUser.ID)
	if err != nil {

		// If there is an error in setting the cache, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
	http.SetCookie(w, &http.Cookie{
		Name:    tokenName,
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}
