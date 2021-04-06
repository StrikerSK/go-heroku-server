package api

import (
	"encoding/json"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go-heroku-server/api/user"
	"log"
	"net/http"

	"go-heroku-server/api/utils"
	"go-heroku-server/config"
)

func addUser(userToAdd user.User) {
	config.DBConnection.NewRecord(userToAdd)
	config.DBConnection.Create(&userToAdd)

	log.Println("Inserted username: " + userToAdd.Username + ".")
}

func GetUserList(w http.ResponseWriter, r *http.Request) {

	var users []user.User
	config.DBConnection.Find(&users)
	for index, user := range users {
		config.DBConnection.Model(&user).Related(&user.Address, "Address")
		users[index] = user
	}

	json.NewEncoder(w).Encode(users)
	log.Println("Retrieved list of users")

}

func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var user user.User
	err := decoder.Decode(&user)

	if err != nil {
		panic(err)
	}

	if verifyNoUserExists(user.Username) {
		addUser(user)
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
}

func EditUser(w http.ResponseWriter, r *http.Request) {

	var name user.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&name)
	if err != nil {
		panic(err)
	}

	if params := r.FormValue("id"); params != "" {

		var newUser user.User

		config.DBConnection.First(&newUser, params)
		config.DBConnection.Model(&newUser).Update(user.User{Username: name.Username, FirstName: name.FirstName, LastName: name.LastName})

		log.Println("User changed")

	} else {

		var newUser user.User
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&name)
		addUser(newUser)

	}
}

func getUserFromDB(username string) (user user.User, userExist bool) {
	userExist = config.DBConnection.Where("username = ?", username).First(&user).RecordNotFound()
	return user, userExist
}

//Function verifies user if it exists and has valid login credentials
func LoginUser(w http.ResponseWriter, r *http.Request) {

	var loginCredentials user.Credentials
	var serverToken user.Token

	json.NewDecoder(r.Body).Decode(&loginCredentials)

	user, userExists := getUserFromDB(loginCredentials.Username)

	log.Print("User is logging \"" + user.Username + "\"")

	if !userExists {

		if loginCredentials.Password == user.Password {

			serverToken = utils.CreateToken(user)
			w.Header().Set("Content-Type", "application/json")
			payload, _ := json.Marshal(serverToken)
			log.Print("User \"" + user.Username + "\" logged successfully!")
			w.Write(payload)

		} else {

			log.Print("User \"" + user.Username + "\" not logged succesfully!")
			http.Error(w, "Unvalid password", 401)

		}

	} else {

		http.Error(w, "User "+loginCredentials.Username+" not found!", 404)

	}

}

//Function verifies user if it exists and has valid login credentials
func GetUserDetail(w http.ResponseWriter, r *http.Request) {

	var requestedUser user.User
	receivedToken := (r.Header.Get("Authorization"))
	userId, _ := utils.GetIdFromToken(receivedToken)
	config.DBConnection.Where("id = ?", userId).Find(&requestedUser).Find(&requestedUser.Address)
	w.Header().Set("Content-Type", "application/json")
	payload, _ := json.Marshal(requestedUser)
	w.Write(payload)
}

//Function verifies user can be registered to database
func verifyNoUserExists(username string) (userExists bool) {

	var requestedUser user.User
	userExists = config.DBConnection.Where("username = ?", username).Find(&requestedUser).RecordNotFound()
	return userExists

}
