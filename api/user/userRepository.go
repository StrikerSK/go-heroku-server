package user

import "go-heroku-server/config"

func createUser(newUser User) {
	config.DBConnection.NewRecord(newUser)
	config.DBConnection.Create(&newUser)
}

//Function retrieves user and flag if exists can be registered to database
func getUserFromDB(username string) (user User, userExist bool) {
	userExist = config.DBConnection.Where("username = ?", username).First(&user).RecordNotFound()
	return user, !userExist
}
