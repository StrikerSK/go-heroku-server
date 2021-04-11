package user

import (
	"encoding/json"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go-heroku-server/api/types"
	"log"
	"net/http"

	"go-heroku-server/config"
)

func getUserList(w http.ResponseWriter, r *http.Request) {

	var users []User
	config.DBConnection.Find(&users)
	for index, currentUser := range users {
		config.DBConnection.Model(&currentUser).Related(&currentUser.Address, "Address")
		users[index] = currentUser
	}

	_ = json.NewEncoder(w).Encode(users)
	log.Println("Retrieved list of users")

}

func registerUser(w http.ResponseWriter, r *http.Request) {
	userBody := r.Context().Value(userBodyContextKey).(User)

	if err := addUser(userBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func addUser(userBody User) (err error) {
	if _, err = getUserByUsername(userBody.Username); err != nil {
		//userBody.decryptPassword()
		userBody.setRole()
		createUser(userBody)
		log.Print("User has been created")
		return
	} else {
		log.Print("User exists in database")
		return
	}
}

func editUser(w http.ResponseWriter, r *http.Request) {

	var updatedUser User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("Decoding issues for user")
		return
	}

	userID := r.Context().Value(UserIdContextKey).(uint)
	persistedUser, err := getUserByID(userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	persistedUser.ID = userID
	if err = updateUser(persistedUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("Something went wrong during user update")
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	claimId := r.Context().Value(UserIdContextKey)
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
