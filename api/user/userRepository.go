package user

import "go-heroku-server/config"

func readUser(username string) (user User) {
	config.DBConnection.Where("username = ?", username).Find(&user)
	return
}

func createUser(newUser User) {
	config.DBConnection.NewRecord(newUser)
	config.DBConnection.Create(&newUser)
}

//Function verifies user can be registered to database
func userExists(username string) (userExists bool) {

	var requestedUser User
	userExists = config.DBConnection.Where("username = ?", username).Find(&requestedUser).RecordNotFound()
	return

}
