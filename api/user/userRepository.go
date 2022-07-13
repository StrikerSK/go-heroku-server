package user

import (
	userDomains "go-heroku-server/api/user/domain"
	"go-heroku-server/config"
)

func createUser(newUser userDomains.User) {
	instance := config.GetDatabaseInstance()
	instance.NewRecord(newUser)
	instance.Create(&newUser)
}

//Function retrieves user and flag if exists can be registered to database
func readUsersFromRepository() (user []userDomains.User, err error) {
	err = config.GetDatabaseInstance().Preload("Address").Find(&user).Error
	return
}

//Function retrieves user and flag if exists can be registered to database
func getUserByID(userID interface{}) (user userDomains.User, err error) {
	err = config.GetDatabaseInstance().Preload("Address").Where("id = ?", userID).First(&user).Error
	return
}

//Function retrieves user and flag if exists can be registered to database
func getUserByUsername(username string) (user userDomains.User, err error) {
	err = config.GetDatabaseInstance().Preload("Address").Where("username = ?", username).First(&user).Error
	return
}

func updateUser(updatedUser userDomains.User) (err error) {
	err = config.GetDatabaseInstance().Model(&userDomains.User{}).Update(&updatedUser).Error
	return
}
