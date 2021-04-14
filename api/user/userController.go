package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go-heroku-server/api/src"
	"go-heroku-server/config"
	"log"
	"net/http"
	"os/user"
)

const (
	userIdContextKey   = "user_id"
	userBodyContextKey = "user_body"
)

func EnrichRouterWithUser(router *mux.Router) {

	config.DBConnection.AutoMigrate(&user.User{})

	userSubroute := router.PathPrefix("/user").Subrouter()
	userSubroute.HandleFunc("/login", controllerLogin).Methods("POST")
	userSubroute.Handle("/register", resolveUser(http.HandlerFunc(controllerRegisterUser))).Methods("POST")

	userSubroute.Handle("/", verifyCookieSession(resolveUser(http.HandlerFunc(controllerEditUser)))).Methods("PUT")
	userSubroute.Handle("/", verifyCookieSession(http.HandlerFunc(controllerGetUser))).Methods("GET")

	jwtSubroute := router.PathPrefix("/jwt").Subrouter()
	jwtSubroute.HandleFunc("/login", LoginUser).Methods("POST")

	jwtSubroute.Handle("/", VerifyJwtToken(resolveUser(http.HandlerFunc(controllerEditUser)))).Methods("PUT")
	jwtSubroute.Handle("/", VerifyJwtToken(http.HandlerFunc(controllerGetUser))).Methods("GET")

	usersSubroute := router.PathPrefix("/users").Subrouter()
	usersSubroute.Handle("/", verifyCookieSession(http.HandlerFunc(controllerGetUserList))).Methods("GET")
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

		ctx := context.WithValue(r.Context(), userIdContextKey, response)
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
		ctx := context.WithValue(r.Context(), userIdContextKey, userClaim.Id)
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

func ResolveUserContext(context context.Context) (uint, src.IResponse) {
	value, ok := context.Value(userIdContextKey).(uint)
	if !ok {
		return 0, src.RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("cannot resolve userID from context"),
		}
	}

	return value, nil
}

// LoginUser godoc
// @Summary Login user
// @Description Login user
// @Tags user
// @Accept json
// @Produce  json
// @Param Credentials body user.Credentials true "Credentials of logged user"
// @Success 200 {object} user.Token "Token generate for user's session"
// @Failure 401 {none} string "Credentials are not valid"
// @Failure 404 {string} string "User does not exists"
// @Router /user [post]
func controllerLogin(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	// Get the JSON body and decode into credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err.Error())
		return
	}

	cookies, outputError := login(credentials)
	if outputError != nil {
		w.WriteHeader(outputError.StatusCode)
		log.Print(outputError.Err)
		return
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
	http.SetCookie(w, cookies)
}

// GetUsers godoc
// @Summary Read all users
// @Description Retrieving all available users
// @Tags user
// @Accept json
// @Produce  json
// @Param Authorization header string true "JWT Token provided from login"
// @Success 200 {array} user.User "List of all stored users in application"
// @Failure 401 {string} string "User is not authorized"
// @Router /users [get]
func controllerGetUserList(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(getUserList())
	log.Println("Retrieved list of users")
}

// RegisterUser godoc
// @Summary Create new user
// @Description Endpoint to provide registration option
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body User true "Body of registered user"
// @Success 201 {object} string "User has been store to database successfully"
// @Failure 409 {object} src.ResponseObject "User with provided username already exists in database"
// @Router /user/register [post]
func controllerRegisterUser(w http.ResponseWriter, r *http.Request) {
	userBody := r.Context().Value(userBodyContextKey).(User)

	if requestError := addUser(userBody); requestError != nil {
		w.WriteHeader(requestError.StatusCode)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func controllerEditUser(w http.ResponseWriter, r *http.Request) {
	userBody := r.Context().Value(userBodyContextKey).(User)
	if requestError := editUser(userBody); requestError != nil {
		w.WriteHeader(requestError.StatusCode)
		log.Print(requestError.Err)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}

// GetUserDetail godoc
// @Summary Get user details
// @Description Detail of currently logged user
// @Tags user
// @Produce  json
// @Param Authorization header string true "JWT Token provided from login"
// @Success 200 {object} user.User "Details of currently logged user"
// @Failure 401 {none} string "User is not authorized"
// @Failure 404 {string} string "User does not exists"
// @Router /user [get]
func controllerGetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIdContextKey)
	userEntity, requestError := getUser(userID)
	if requestError != nil {
		w.WriteHeader(requestError.StatusCode)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		payload, _ := json.Marshal(userEntity)
		_, _ = w.Write(payload)
	}
}
