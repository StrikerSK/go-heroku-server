package user

import (
	"errors"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go-heroku-server/api/src/responses"
	"go-heroku-server/api/types"
	"log"
	"net/http"
)

func getUserList() responses.IResponse {
	if users, err := readUsersFromRepository(); err != nil {
		log.Printf("Users read: %s\n", err.Error())
		return responses.NewErrorResponse(http.StatusBadRequest, err)
	} else {
		for i := range users {
			users[i].Password = ""
		}
		return responses.NewResponse(users)
	}
}

func addUser(userBody User) responses.IResponse {
	if _, err := getUserByUsername(userBody.Username); err != nil {
		//userBody.decryptPassword()
		userBody.setRole()
		createUser(userBody)
		log.Printf("User add: created\n")
		return responses.NewEmptyResponse(http.StatusCreated)
	} else {
		log.Printf("User add: already exists\n")
		return responses.NewEmptyResponse(http.StatusConflict)
	}
}

func editUser(updatedUser User) responses.IResponse {
	persistedUser, err := getUserByID(updatedUser.ID)
	if err != nil {
		log.Printf("User [%d] edit: %s\n", updatedUser.ID, err.Error())
		return responses.NewEmptyResponse(http.StatusNotFound)
	}

	persistedUser.ID = updatedUser.ID
	if err = updateUser(persistedUser); err != nil {
		log.Printf("User [%d] edit: %s", updatedUser.ID, err.Error())
		return responses.NewErrorResponse(http.StatusBadRequest, errors.New("user update failed"))
	}

	return responses.NewEmptyResponse(http.StatusOK)
}

func getUser(userID interface{}) responses.IResponse {
	if requestedUser, err := getUserByID(userID); err != nil {
		log.Printf("User [%d] read: %s", userID, err.Error())
		return responses.NewErrorResponse(http.StatusNotFound, err)
	} else {
		requestedUser.Password = ""
		return responses.NewResponse(requestedUser)
	}
}

func InitAdminUser() {
	user := User{
		FirstName: "admin",
		LastName:  "admin",
		Role:      AdminRole,
		Address: types.Address{
			Street: "Admin",
			City:   "Admin",
			Zip:    "Admin",
		},
		Credentials: Credentials{
			Username: "admin",
			Password: "admin",
		},
	}

	//user.decryptPassword()

	_ = addUser(user)
}

func InitCommonUser() {
	user := User{
		FirstName: "tester",
		LastName:  "tester",
		Role:      UserRole,
		Address: types.Address{
			Street: "Tester",
			City:   "Tester",
			Zip:    "Tester",
		},
		Credentials: Credentials{
			Username: "admin",
			Password: "admin",
		},
	}

	//user.decryptPassword()

	_ = addUser(user)
}
