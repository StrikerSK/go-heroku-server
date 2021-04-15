package user

import (
	"errors"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go-heroku-server/api/src"
	"go-heroku-server/api/types"
	"log"
	"net/http"

	"go-heroku-server/config"
)

func getUserList() (users []User) {
	config.DBConnection.Find(&users)
	for index, currentUser := range users {
		config.DBConnection.Model(&currentUser).Related(&currentUser.Address, "Address")
		users[index] = currentUser
	}
	return
}

func addUser(userBody User) *src.RequestError {
	var requestError src.RequestError

	if _, err := getUserByUsername(userBody.Username); err != nil {
		userBody.decryptPassword()
		userBody.setRole()
		createUser(userBody)
		log.Print("User has been created")
		return nil
	} else {
		log.Print("User exists in database")
		requestError.Err = errors.New("user exists in database")
		requestError.StatusCode = http.StatusConflict
		return &requestError
	}
}

func editUser(updatedUser User) *src.RequestError {
	var requestError src.RequestError

	persistedUser, err := getUserByID(updatedUser.ID)
	if err != nil {
		requestError.Err = errors.New("user not found")
		requestError.StatusCode = http.StatusNotFound
		return &requestError
	}

	persistedUser.ID = updatedUser.ID
	if err = updateUser(persistedUser); err != nil {
		requestError.Err = errors.New("user update failed")
		requestError.StatusCode = http.StatusBadRequest
		return &requestError
	}

	return nil
}

func getUser(userID interface{}) (*User, *src.RequestError) {
	var requestError src.RequestError
	requestedUser, err := getUserByID(userID)
	if err != nil {
		requestError.Err = err
		requestError.StatusCode = http.StatusNotFound
		log.Print(err)
		return nil, &requestError
	}
	return &requestedUser, nil
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

	user.decryptPassword()

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

	user.decryptPassword()

	_ = addUser(user)
}
