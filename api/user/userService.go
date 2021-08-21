package user

import (
	"encoding/json"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go-heroku-server/api/src/responses"
	customAuth "go-heroku-server/api/user/auth"
	"log"
	"net/http"
)

func getUserList() responses.IResponse {
	if users, err := readUsersFromRepository(); err != nil {
		log.Printf("Users read: %s\n", err.Error())
		return responses.CreateResponse(http.StatusBadRequest, nil)
	} else {
		for i := range users {
			users[i].clearPassword()
		}
		return responses.CreateResponse(http.StatusOK, users)
	}
}

func addUser(userBody User) responses.IResponse {
	if _, err := getUserByUsername(userBody.Username); err != nil {
		//userBody.decryptPassword()
		userBody.setRole()
		createUser(userBody)
		log.Printf("User add: created\n")
		return responses.CreateResponse(http.StatusCreated, nil)
	} else {
		log.Printf("User add: already exists\n")
		return responses.CreateResponse(http.StatusConflict, nil)
	}
}

func editUser(updatedUser User) responses.IResponse {
	persistedUser, err := getUserByID(updatedUser.ID)
	if err != nil {
		log.Printf("User [%d] edit: %s\n", updatedUser.ID, err.Error())
		return responses.CreateResponse(http.StatusNotFound, nil)
	}

	persistedUser.ID = updatedUser.ID
	if err = updateUser(persistedUser); err != nil {
		log.Printf("User [%d] edit: %s", updatedUser.ID, err.Error())
		return responses.CreateResponse(http.StatusBadRequest, nil)
	}

	return responses.CreateResponse(http.StatusOK, nil)
}

func getUser(userID interface{}) responses.IResponse {
	if requestedUser, err := getUserByID(userID); err != nil {
		log.Printf("User [%d] read: %s", userID, err.Error())
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
		log.Printf("Logging error: %s\n", err)
		responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
		return
	}

	persistedUser, err := getUserByUsername(credentials.Username)
	if err != nil {
		log.Printf("Logging error: %s\n", err)
		responses.CreateResponse(http.StatusUnauthorized, nil).WriteResponse(w)
		return
	}

	//if err = persistedUser.validatePassword(credentials.Password); err != nil {
	if !persistedUser.validatePassword(credentials.Password) {
		responses.CreateResponse(http.StatusUnauthorized, nil).WriteResponse(w)
		return
	}

	signetToken, err := customAuth.CreateToken(persistedUser)
	if err != nil {
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	log.Printf("User [%d] login: success\n", persistedUser.ID)

	res := responses.CreateResponse(http.StatusOK, &Token{Token: signetToken})
	res.WriteResponse(w)
	return
}
