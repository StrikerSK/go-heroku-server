package userHandlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"go-heroku-server/api/src/responses"
	customAuth "go-heroku-server/api/user/auth"
	userDomains "go-heroku-server/api/user/domain"
	userPorts "go-heroku-server/api/user/ports"
	"log"
	"net/http"
)

const (
	userIdContextKey   = "user_id"
	userBodyContextKey = "user_body"
)

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
	userRoute.Handle("/register", h.middleware.ResolveUser(http.HandlerFunc(h.createUser))).Methods(http.MethodPost)
	userRoute.HandleFunc("/login", h.login).Methods(http.MethodPost)
	userRoute.Handle("/", h.middleware.VerifyToken(h.middleware.ResolveUser(http.HandlerFunc(h.updateUser)))).Methods(http.MethodPut)
	userRoute.Handle("/", h.middleware.VerifyToken(http.HandlerFunc(h.readUser))).Methods(http.MethodGet)

	usersRoute := router.PathPrefix("/users").Subrouter()
	usersRoute.Handle("/", h.middleware.VerifyToken(http.HandlerFunc(h.readUsers))).Methods(http.MethodGet)
}

func ResolveUserContext(context context.Context) (uint, responses.IResponse) {
	value, ok := context.Value(userIdContextKey).(uint)
	if !ok {
		log.Println("UserID resolve: cannot resolve from context")
		return 0, responses.CreateResponse(http.StatusInternalServerError, nil)
	}

	return value, nil
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
	user := r.Context().Value(userBodyContextKey).(userDomains.User)
	if _, err := h.userService.ReadUser(user.Username); err != nil {
		if err == gorm.ErrRecordNotFound {
			user.SetRole()
			_ = h.userService.CreateUser(user)
			log.Printf("User add: created\n")
			responses.CreateResponse(http.StatusCreated, nil).WriteResponse(w)
			return
		} else {
			log.Printf("User add: bad request occured\n")
			responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
			return
		}
	} else {
		log.Printf("User add: user already created\n")
		responses.CreateResponse(http.StatusConflict, nil).WriteResponse(w)
		return
	}
}

func (h UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	userBody := r.Context().Value(userBodyContextKey).(userDomains.User)
	if err := h.userService.UpdateUser(userBody); err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("User [%s] edit: user not found\n", userBody.Username)
			responses.CreateResponse(http.StatusNotFound, nil).WriteResponse(w)
			return
		} else {
			log.Printf("User [%s] edit: %v\n", userBody.Username, err.Error())
			responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
			return
		}
	} else {
		responses.CreateResponse(http.StatusOK, nil).WriteResponse(w)
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

	var credentials userDomains.Credentials

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
