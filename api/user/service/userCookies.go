package userServices

import (
	uuid "github.com/satori/go.uuid"
	"go-heroku-server/api/src/responses"
	userDomains "go-heroku-server/api/user/domain"
	userPorts "go-heroku-server/api/user/ports"
	"go-heroku-server/config"
	"log"
	"net/http"
	"time"
)

const tokenName = "session_token"

type CookieService struct {
	repository userPorts.IUserRepository
}

func NewCookieService(repository userPorts.IUserRepository) CookieService {
	return CookieService{
		repository: repository,
	}
}

func (s CookieService) Login(credentials userDomains.UserCredentials) (http.Cookie, responses.IResponse) {
	var requestError responses.IResponse

	persistedUser, err := s.repository.ReadUserByUsername(credentials.Username)
	if err != nil {
		requestError = responses.CreateResponse(http.StatusBadRequest, err)
		return http.Cookie{}, requestError
	}

	// If a password exists for the given user
	// AND, if it is the same as the password we received, then we can move ahead
	// if NOT, then we return an "Unauthorized" status
	//if err = persistedUser.validatePassword(credentials.Password); err != nil {
	if !persistedUser.ValidatePassword(credentials.Password) {
		log.Println("Login with Cookies: password are not matching")
		requestError = responses.CreateResponse(http.StatusUnauthorized, nil)
		return http.Cookie{}, requestError
	}

	// Create a new random session token
	sessionToken := uuid.NewV4().String()
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	_, err = config.GetCacheInstance().Do("SETEX", sessionToken, "120", persistedUser.Username)
	if err != nil {
		// If there is an error in setting the cache, return an internal server error
		log.Printf("Login with Cookies: %s\n", err.Error())
		requestError = responses.CreateResponse(http.StatusInternalServerError, nil)
		return http.Cookie{}, requestError
	}

	cookie := http.Cookie{
		Name:    tokenName,
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	}

	return cookie, nil
}
