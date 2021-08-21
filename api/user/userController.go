package user

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-heroku-server/api/src/responses"
	customAuth "go-heroku-server/api/user/auth"
	"go-heroku-server/config"
	"log"
	"net/http"
)

const (
	userIdContextKey   = "user_id"
	userBodyContextKey = "user_body"
)

func EnrichRouterWithUser(router *mux.Router) {

	userSubroute := router.PathPrefix("/user").Subrouter()
	userSubroute.HandleFunc("/login", loginGeneratingCookie).Methods(http.MethodPost)
	userSubroute.Handle("/register", resolveUser(http.HandlerFunc(controllerRegisterUser))).Methods(http.MethodPost)

	userSubroute.Handle("/", verifyCookieSession(resolveUser(http.HandlerFunc(controllerEditUser)))).Methods(http.MethodPut)
	userSubroute.Handle("/", verifyCookieSession(http.HandlerFunc(controllerGetUser))).Methods(http.MethodGet)

	jwtSubroute := router.PathPrefix("/jwt").Subrouter()
	jwtSubroute.HandleFunc("/login", LoginUser).Methods(http.MethodPost)

	jwtSubroute.Handle("/", VerifyJwtToken(resolveUser(http.HandlerFunc(controllerEditUser)))).Methods(http.MethodPut)
	jwtSubroute.Handle("/", VerifyJwtToken(http.HandlerFunc(controllerGetUser))).Methods(http.MethodGet)

	usersSubroute := router.PathPrefix("/users").Subrouter()
	usersSubroute.Handle("/", VerifyJwtToken(http.HandlerFunc(controllerGetUserList))).Methods(http.MethodGet)
}

func verifyCookieSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(tokenName)
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

func VerifyJwtToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Println("Verify JWT Token: Cannot find header")
			responses.CreateResponse(http.StatusUnauthorized, nil).WriteResponse(w)
			return
		}
		userClaim, err := customAuth.ParseToken(token)
		if err != nil {
			responses.CreateResponse(http.StatusUnauthorized, nil).WriteResponse(w)
			return
		}
		ctx := context.WithValue(r.Context(), userIdContextKey, userClaim.Id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func resolveUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user User
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&user); err != nil {
			log.Printf("Resolve user: %s\n", err.Error())
			responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
			return
		}

		ctx := context.WithValue(r.Context(), userBodyContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ResolveUserContext(context context.Context) (uint, responses.IResponse) {
	value, ok := context.Value(userIdContextKey).(uint)
	if !ok {
		log.Println("UserID resolve: cannot resolve from context")
		return 0, responses.CreateResponse(http.StatusInternalServerError, nil)
	}

	return value, nil
}

func loginGeneratingCookie(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	// Get the JSON body and decode into credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Printf("Cookie Login: %s\n", err.Error())
		responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
		return
	}

	cookies, outputError := loginWithCookies(credentials)
	if outputError != nil {
		outputError.WriteResponse(w)
		return
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
	http.SetCookie(w, &cookies)
}

func controllerGetUserList(w http.ResponseWriter, r *http.Request) {
	getUserList().WriteResponse(w)
}

func controllerRegisterUser(w http.ResponseWriter, r *http.Request) {
	userBody := r.Context().Value(userBodyContextKey).(User)
	addUser(userBody).WriteResponse(w)
}

func controllerEditUser(w http.ResponseWriter, r *http.Request) {
	userBody := r.Context().Value(userBodyContextKey).(User)
	editUser(userBody).WriteResponse(w)
}

func controllerGetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIdContextKey)
	getUser(userID).WriteResponse(w)
}
