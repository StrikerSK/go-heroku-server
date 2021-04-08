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

func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var newUser User
	err := decoder.Decode(&newUser)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if userExists(newUser.Username) {
		createUser(newUser)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func EditUser(w http.ResponseWriter, r *http.Request) {

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

func getUserFromDB(username string) (user User, userExist bool) {
	userExist = config.DBConnection.Where("username = ?", username).First(&user).RecordNotFound()
	return user, !userExist
}

//Function verifies user if it exists and has valid login credentials
func GetUserDetail(w http.ResponseWriter, r *http.Request) {
	userId := receiveCookie(r)
	getUser(userId, w)
}

func GetUserDetailFromJWT(w http.ResponseWriter, r *http.Request) {
	claims, err := ParseToken(r.Header.Get("Authorization"))
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	getUser(claims.Id, w)
}

func getUser(id interface{}, w http.ResponseWriter) {
	var requestedUser User
	config.DBConnection.Where("id = ?", id).Find(&requestedUser).Find(&requestedUser.Address)
	w.Header().Set("Content-Type", "application/json")
	payload, _ := json.Marshal(requestedUser)
	_, _ = w.Write(payload)
}
