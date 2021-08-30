package user

import (
	"encoding/json"
	"errors"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	customAuth "go-heroku-server/src/api/user/auth"
	"go-heroku-server/src/responses"
	"log"
	"net/http"
)

func getUserList() responses.IResponse {
	users, err := readUsersFromRepository()
	if err != nil {
		log.Printf("Users read: %v\n", err)
		return responses.CreateResponse(http.StatusBadRequest, nil)
	}

	for i := range users {
		users[i].clearPassword()
	}

	return responses.CreateResponse(http.StatusOK, users)
}

func addUser(userBody User) responses.IResponse {
	if _, err := getUserByUsername(userBody.Username); err != nil {
		//userBody.decryptPassword()
		userBody.setRole()
		createUser(userBody)
		log.Printf("User create: success\n")
		return responses.CreateResponse(http.StatusCreated, nil)
	} else {
		log.Printf("User create: conflict\n")
		return responses.CreateResponse(http.StatusConflict, nil)
	}
}

func editUser(updatedUser User) responses.IResponse {
	persistedUser, err := getUserByID(updatedUser.ID)
	if err != nil {
		log.Printf("User [%d] update: %v\n", updatedUser.ID, err)
		return responses.CreateResponse(http.StatusNotFound, nil)
	}

	persistedUser.ID = updatedUser.ID
	if err = updateUser(persistedUser); err != nil {
		log.Printf("User [%d] update: %v\n", updatedUser.ID, err)
		return responses.CreateResponse(http.StatusBadRequest, nil)
	}

	return responses.CreateResponse(http.StatusOK, nil)
}

func getUser(userID interface{}) responses.IResponse {
	if requestedUser, err := getUserByID(userID); err != nil {
		log.Printf("User [%d] read: %v\n", userID, err)
		return responses.CreateResponse(http.StatusNotFound, nil)
	} else {
		requestedUser.clearPassword()
		return responses.CreateResponse(http.StatusOK, requestedUser)
	}
}

//Function verifies user if it exists and has valid login credentials
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		log.Printf("User login: %v\n", err)
		responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
		return
	}

	persistedUser, err := validateUser(credentials)
	if err != nil {
		log.Printf("User login: %v\n", err)
		responses.CreateResponse(http.StatusUnauthorized, nil).WriteResponse(w)
		return
	}

	signetToken, err := customAuth.CreateToken(persistedUser)
	if err != nil {
		log.Printf("User login: %v\n", err)
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	log.Printf("User [%d] login: success\n", persistedUser.ID)
	responses.CreateResponse(http.StatusOK, &Token{Token: signetToken}).WriteResponse(w)
	return
}

func validateUser(credentials Credentials) (user User, err error) {
	user, err = getUserByUsername(credentials.Username)
	if err != nil {
		return
	}

	// If a password exists for the given user AND passwords are matching
	if !user.validatePassword(credentials.Password) {
		return User{}, errors.New("password are not matching")
	}

	return
}
