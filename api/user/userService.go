package user

import (
	"errors"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go-heroku-server/api/src"
	"go-heroku-server/api/types"
	"go-heroku-server/config"
	"log"
	"net/http"
)

func getUserList() (users []User) {
	config.DBConnection.Find(&users)
	for index, currentUser := range users {
		config.DBConnection.Model(&currentUser).Related(&currentUser.Address, "Address")
		users[index] = currentUser
	}
	return
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

func getUser(userID interface{}) (res src.IResponse) {
	if requestedUser, err := getUserByID(userID); err != nil {
		res = src.NewErrorResponse(http.StatusNotFound, err)
	} else {
		res = src.NewResponse(requestedUser)
	}
	return
}

func InitAdminUser() {
	user := User{
		Username:  "admin",
		Password:  "admin",
		FirstName: "admin",
		LastName:  "admin",
		Role:      AdminRole,
		Address: types.Address{
			Street: "Admin",
			City:   "Admin",
			Zip:    "Admin",
		},
	}

	//user.decryptPassword()

	_ = addUser(user)
}

func InitCommonUser() {
	user := User{
		Username:  "tester",
		Password:  "tester",
		FirstName: "tester",
		LastName:  "tester",
		Role:      UserRole,
		Address: types.Address{
			Street: "Tester",
			City:   "Tester",
			Zip:    "Tester",
		},
	}

	//user.decryptPassword()

	_ = addUser(user)
}
