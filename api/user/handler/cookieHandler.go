package userHandlers

import (
	"context"
	"encoding/json"
	"go-heroku-server/api/src/responses"
	userDomains "go-heroku-server/api/user/domain"
	userServices "go-heroku-server/api/user/service"
	"go-heroku-server/config"
	"log"
	"net/http"
)

type UserCookiesHandler struct {
	service userServices.CookieService
}

func NewUserCookiesHandler(service userServices.CookieService) UserCookiesHandler {
	return UserCookiesHandler{
		service: service,
	}
}

func (h UserCookiesHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials userDomains.Credentials
	// Get the JSON body and decode into credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Printf("Cookie Login: %s\n", err.Error())
		responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
		return
	}

	cookies, outputError := h.service.Login(credentials)
	if outputError != nil {
		outputError.WriteResponse(w)
		return
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
	http.SetCookie(w, &cookies)
}

func (h UserCookiesHandler) VerifySession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				log.Printf("Verify cookie: %s\n", err.Error())
				responses.CreateResponse(http.StatusUnauthorized, nil).WriteResponse(w)
				return
			}
			// For any other type of error, return a bad request status
			log.Printf("Verify cookie: %s\n", err.Error())
			responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
			return
		}

		sessionToken := c.Value

		// We then get the name of the user from our cache, where we set the session token
		response, err := config.GetCacheInstance().Do("GET", sessionToken)
		if err != nil {
			log.Printf("Cache login: %s\n", err.Error())
			responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
			return
		}
		if response == nil {
			responses.CreateResponse(http.StatusUnauthorized, nil).WriteResponse(w)
			return
		}

		ctx := context.WithValue(r.Context(), userIdContextKey, response)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
