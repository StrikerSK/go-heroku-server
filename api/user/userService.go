package user

import (
	"encoding/json"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"

	"go-heroku-server/config"
)

func GetUserList(w http.ResponseWriter, r *http.Request) {

	var users []User
	config.DBConnection.Find(&users)
	for index, currentUser := range users {
		config.DBConnection.Model(&currentUser).Related(&currentUser.Address, "Address")
		users[index] = currentUser
	}

	_ = json.NewEncoder(w).Encode(users)
	log.Println("Retrieved list of users")

}

func registerNewUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var newUser User
	err := decoder.Decode(&newUser)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	_, userExists := getUserFromDB(newUser.Username)
	if userExists {
		//encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		//newUser.Password = string(encryptedPassword)
		createUser(newUser)
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func editUser(w http.ResponseWriter, r *http.Request) {

	var updatedUser User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		panic(err)
	}

	if params := r.FormValue("id"); params != "" {

		var persistedUser User

		config.DBConnection.First(&persistedUser, params)
		config.DBConnection.Model(&persistedUser).Update(User{Username: updatedUser.Username, FirstName: updatedUser.FirstName, LastName: updatedUser.LastName})

		log.Println("User changed")

	} else {

		var newUser User
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&updatedUser)
		createUser(newUser)

	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	claimId := r.Context().Value("user_id")
	if claimId == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var requestedUser User
	config.DBConnection.Where("id = ?", claimId).Find(&requestedUser).Find(&requestedUser.Address)
	w.Header().Set("Content-Type", "application/json")
	payload, _ := json.Marshal(requestedUser)
	_, _ = w.Write(payload)
}
