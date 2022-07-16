package userHandlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-heroku-server/api/src/responses"
	customAuth "go-heroku-server/api/user/auth"
	userDomains "go-heroku-server/api/user/domain"
	userPorts "go-heroku-server/api/user/ports"
	"io"
	"log"
	"net/http"
)

const userIdContextKey = "user_id"

type UserHandler struct {
	userService  userPorts.IUserService
	middleware   UserAuthMiddleware
	tokenService customAuth.TokenService
}

func NewUserHandler(userService userPorts.IUserService, middleware UserAuthMiddleware, tokenService customAuth.TokenService) UserHandler {
	return UserHandler{
		userService:  userService,
		middleware:   middleware,
		tokenService: tokenService,
	}
}

func (h UserHandler) EnrichRouter(router *mux.Router) {
	userRoute := router.PathPrefix("/user").Subrouter()
	userRoute.Handle("/register", http.HandlerFunc(h.createUser)).Methods(http.MethodPost)
	userRoute.Handle("/login", http.HandlerFunc(h.login)).Methods(http.MethodPost)

	userRoute.Handle("", h.middleware.VerifyToken(http.HandlerFunc(h.updateUser))).Methods(http.MethodPut)
	userRoute.Handle("", h.middleware.VerifyToken(http.HandlerFunc(h.readUser))).Methods(http.MethodGet)

	usersRoute := router.PathPrefix("/users").Subrouter()
	usersRoute.Handle("", h.middleware.VerifyToken(http.HandlerFunc(h.readUsers))).Methods(http.MethodGet)
}

func (h UserHandler) readUsers(w http.ResponseWriter, r *http.Request) {
	if users, err := h.userService.ReadUsers(); err != nil {
		log.Printf("Users read: %s\n", err.Error())
		responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
		return
	} else {
		for i := range users {
			users[i].ClearPassword()
		}
		responses.CreateResponse(http.StatusOK, users).WriteResponse(w)
		return
	}
}

func (h UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.parseUser(r.Body)
	if err != nil {
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	err = h.userService.CreateUser(user)
	if err != nil {
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	} else {
		responses.CreateResponse(http.StatusCreated, nil).WriteResponse(w)
		return
	}
}

func (h UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	userBody, err := h.parseUser(r.Body)
	if err != nil {
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	username, _ := h.middleware.GetUserFromContext(r.Context())
	if username != userBody.Username {
		log.Printf("user cannot update without correct token")
		responses.CreateResponse(http.StatusForbidden, nil).WriteResponse(w)
		return
	}

	err = h.userService.UpdateUser(userBody)
	if err != nil {
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	} else {
		responses.CreateResponse(http.StatusOK, nil).WriteResponse(w)
		return
	}
}

func (h UserHandler) readUser(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userIdContextKey).(string)
	if requestedUser, err := h.userService.ReadUser(username); err != nil {
		log.Printf("User [%s] read: %s", username, err.Error())
		responses.CreateResponse(http.StatusNotFound, nil).WriteResponse(w)
		return
	} else {
		requestedUser.ClearPassword()
		responses.CreateResponse(http.StatusOK, requestedUser).WriteResponse(w)
		return
	}
}

//Function verifies user if it exists and has valid login credentials
func (h UserHandler) login(w http.ResponseWriter, r *http.Request) {
	var credentials userDomains.UserCredentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Printf("Logging error: %s\n", err)
		responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
		return
	}

	persistedUser, err := h.userService.ReadUser(credentials.Username)
	if err != nil {
		log.Printf("Logging error: %s\n", err)
		responses.CreateResponse(http.StatusUnauthorized, nil).WriteResponse(w)
		return
	}

	//if err = persistedUser.validatePassword(credentials.Password); err != nil {
	if !persistedUser.ValidatePassword(credentials.Password) {
		responses.CreateResponse(http.StatusUnauthorized, nil).WriteResponse(w)
		return
	}

	signetToken, err := h.tokenService.CreateToken(persistedUser)
	if err != nil {
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	log.Printf("User [%s] login: success\n", persistedUser.Username)

	responses.CreateResponse(http.StatusOK, map[string]string{"token": signetToken}).WriteResponse(w)
	return
}

func (h UserHandler) parseUser(body io.ReadCloser) (userDomains.User, error) {
	var user userDomains.User
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&user); err != nil {
		log.Printf("User parsing error: %v\n", err)
		return userDomains.User{}, err
	}

	return user, nil
}
