package user

import (
	"github.com/gorilla/mux"
	"go-heroku-server/config"
	"net/http"
)

func UserEnrichRouter(router *mux.Router) {

	userSubroute := router.PathPrefix("/user").Subrouter()
	userSubroute.HandleFunc("/login", Login).Methods("POST")
	userSubroute.HandleFunc("/register", RegisterNewUser).Methods("POST")

	userSubroute.Handle("/", verifyCookieSession(http.HandlerFunc(EditUser))).Methods("PUT")
	userSubroute.Handle("/", verifyCookieSession(http.HandlerFunc(GetUserDetail))).Methods("GET")

	jwtSubroute := router.PathPrefix("/jwt").Subrouter()
	jwtSubroute.HandleFunc("/login", LoginUser).Methods("POST")
	jwtSubroute.HandleFunc("/register", RegisterNewUser).Methods("POST")

	jwtSubroute.Handle("/", verifyJwtToken(http.HandlerFunc(EditUser))).Methods("PUT")
	jwtSubroute.Handle("/", verifyJwtToken(http.HandlerFunc(GetUserDetailFromJWT))).Methods("GET")

	usersSubroute := router.PathPrefix("/users").Subrouter()
	usersSubroute.Handle("/", verifyCookieSession(http.HandlerFunc(GetUserList))).Methods("GET")
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

		next.ServeHTTP(w, r)
	})
}

func verifyJwtToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
