package api

import (
	"encoding/json"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"

	"go-heroku-server/api/types"
	"go-heroku-server/api/utils"
	"go-heroku-server/config"
)

var tokenEncodeString = []byte("Wow, much safe")

func addUser(username, password, firstname, lastname string, adress types.Address) {
	db, err := config.CreateDatabase()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	user := types.User{Username: username, Password: password, LastName: lastname, FirstName: firstname, Address: adress}

	db.NewRecord(user)
	db.Create(&user)

	log.Println("Inserted username: " + username + ".")
}

func GetUserList(w http.ResponseWriter, r *http.Request) {

	var users []types.User

	db, err := config.CreateDatabase()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.Find(&users)

	for index, user := range users {

		db.Model(&user).Related(&user.Address, "Address")
		users[index] = user

	}

	json.NewEncoder(w).Encode(users)
	log.Println("Retrieved list of users")

}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var user types.User
	err := decoder.Decode(&user)

	if err != nil {
		panic(err)
	}

	addUser(user.Username, user.Password, user.FirstName, user.LastName, user.Address)
}

func EditName(w http.ResponseWriter, r *http.Request) {

	db, err := config.CreateDatabase()

	if err != nil {
		panic(err)
	}

	var name types.User
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&name)

	if err != nil {
		panic(err)
	}

	if params := r.FormValue("id"); params != "" {

		var newUser types.User

		db.First(&newUser, params)
		db.Model(&newUser).Update(types.User{Username: name.Username, FirstName: name.FirstName, LastName: name.LastName})

		log.Println("User changed")

	} else {

		var name types.User
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&name)
		addUser(name.Username, name.Password, name.FirstName, name.LastName, name.Address)

	}
}

func getUserFromDB(username string) (user types.User, userExist bool) {
	db, err := config.CreateDatabase()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	userExist = db.Where("username = ?", username).First(&user).RecordNotFound()

	return user, userExist
}

//Function verifies user if it exists and has valid login credentials
func LoginUser(w http.ResponseWriter, r *http.Request) {

	var loginCredentials types.UserCredentials
	var serverToken types.Token

	json.NewDecoder(r.Body).Decode(&loginCredentials)

	user, userExists := getUserFromDB(loginCredentials.Username)

	if !userExists {

		if loginCredentials.Password == user.Password {

			serverToken = utils.CreateToken(user)

			w.Header().Set("Content-Type", "application/json")

			payload, _ := json.Marshal(serverToken)

			w.Write([]byte(payload))

		} else {

			http.Error(w, "Unvalid password", 401)

		}

	} else {

		http.Error(w, "User "+loginCredentials.Username+" not found!", 404)

	}

}

//Function verifies user if it exists and has valid login credentials
func GetUserDetail(w http.ResponseWriter, r *http.Request) {

	var requestedUser types.User

	receivedToken := (r.Header.Get("Authorization"))

	userId, _ := utils.GetIdFromToken(receivedToken)

	db, err := config.CreateDatabase()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.Where("id = ?", userId).Find(&requestedUser).Find(&requestedUser.Address)

	w.Header().Set("Content-Type", "application/json")

	payload, _ := json.Marshal(requestedUser)

	w.Write([]byte(payload))

}
