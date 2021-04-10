package user

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-heroku-server/config"
	"log"
	"net/http"
	"os/user"
)

const (
	UserIdContextKey   = "user_id"
	userBodyContextKey = "user_body"
)

func EnrichRouterWithUser(router *mux.Router) {

	config.DBConnection.AutoMigrate(&user.User{})

	userSubroute := router.PathPrefix("/user").Subrouter()
	userSubroute.HandleFunc("/login", login).Methods("POST")
	userSubroute.Handle("/register", resolveUser(http.HandlerFunc(registerUser))).Methods("POST")

	userSubroute.Handle("/", verifyCookieSession(resolveUser(http.HandlerFunc(editUser)))).Methods("PUT")
	userSubroute.Handle("/", verifyCookieSession(http.HandlerFunc(getUser))).Methods("GET")

	jwtSubroute := router.PathPrefix("/jwt").Subrouter()
	jwtSubroute.HandleFunc("/login", LoginUser).Methods("POST")
	jwtSubroute.HandleFunc("/register", registerUser).Methods("POST")

	jwtSubroute.Handle("/", VerifyJwtToken(resolveUser(http.HandlerFunc(editUser)))).Methods("PUT")
	jwtSubroute.Handle("/", VerifyJwtToken(http.HandlerFunc(getUser))).Methods("GET")

	usersSubroute := router.PathPrefix("/users").Subrouter()
	usersSubroute.Handle("/", verifyCookieSession(http.HandlerFunc(getUserList))).Methods("GET")
}

func verifyCookieSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(tokenName)
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sessionToken := c.Value

		// We then get the name of the user from our cache, where we set the session token
		response, err := config.Cache.Do("GET", sessionToken)
		if err != nil {
			// If there is an error fetching from cache, return an internal server error status
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if response == nil {
			// If the session token is not present in cache, return an unauthorized error
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId := receiveCookie(r)
		ctx := context.WithValue(r.Context(), UserIdContextKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func VerifyJwtToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userClaim, err := ParseToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIdContextKey, userClaim.Id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func resolveUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user User
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), userBodyContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
