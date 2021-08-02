package user

import "go-heroku-server/config"

func createUser(newUser User) {
	config.DBConnection.NewRecord(newUser)
	config.DBConnection.Create(&newUser)
}

//Function retrieves user and flag if exists can be registered to database
func readUsersFromRepository() (user []User, err error) {
	err = config.DBConnection.Preload("Address").Find(&user).Error
	return
}

//Function retrieves user and flag if exists can be registered to database
func getUserByID(userID interface{}) (user User, err error) {
	err = config.DBConnection.Preload("Address").Where("id = ?", userID).First(&user).Error
	return
}

//Function retrieves user and flag if exists can be registered to database
func getUserByUsername(username string) (user User, err error) {
	err = config.DBConnection.Preload("Address").Where("username = ?", username).First(&user).Error
	return
}

func updateUser(updatedUser User) (err error) {
	err = config.DBConnection.Model(&User{}).Update(&updatedUser).Error
	return
}
