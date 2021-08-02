package user

import (
	"errors"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go-heroku-server/api/src"
	"go-heroku-server/api/types"
	"log"
	"net/http"
)

func getUserList() src.IResponse {
	if users, err := readUsersFromRepository(); err != nil {
		return src.NewErrorResponse(http.StatusBadRequest, err)
	} else {
		for i := range users {
			users[i].Password = ""
		}
		return src.NewResponse(users)
	}
}

func addUser(userBody User) src.IResponse {
	if _, err := getUserByUsername(userBody.Username); err != nil {
		//userBody.decryptPassword()
		userBody.setRole()
		createUser(userBody)
		log.Print("User has been created")
		return src.NewEmptyResponse(http.StatusCreated)
	} else {
		return src.NewErrorResponse(http.StatusConflict, errors.New("user exists in database"))
	}
}

func editUser(updatedUser User) src.IResponse {
	persistedUser, err := getUserByID(updatedUser.ID)
	if err != nil {
		return src.NewErrorResponse(http.StatusNotFound, errors.New("user not found"))
	}

	persistedUser.ID = updatedUser.ID
	if err = updateUser(persistedUser); err != nil {
		return src.NewErrorResponse(http.StatusBadRequest, errors.New("user update failed"))
	}

	return src.NewEmptyResponse(http.StatusOK)
}

func getUser(userID interface{}) src.IResponse {
	if requestedUser, err := getUserByID(userID); err != nil {
		return src.NewErrorResponse(http.StatusNotFound, err)
	} else {
		requestedUser.Password = ""
		return src.NewResponse(requestedUser)
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
