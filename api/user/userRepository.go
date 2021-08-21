package user

import "go-heroku-server/config"

func createUser(newUser User) {
	config.GetDatabaseInstance().NewRecord(newUser)
	config.GetDatabaseInstance().Create(&newUser)
}

//Function retrieves user and flag if exists can be registered to database
func readUsersFromRepository() (user []User, err error) {
	err = config.GetDatabaseInstance().Preload("Address").Find(&user).Error
	return
}

//Function retrieves user and flag if exists can be registered to database
func getUserByID(userID interface{}) (user User, err error) {
	err = config.GetDatabaseInstance().Preload("Address").Where("id = ?", userID).First(&user).Error
	return
}

//Function retrieves user and flag if exists can be registered to database
func getUserByUsername(username string) (user User, err error) {
	err = config.GetDatabaseInstance().Preload("Address").Where("username = ?", username).First(&user).Error
	return
}

func updateUser(updatedUser User) (err error) {
	err = config.GetDatabaseInstance().Model(&User{}).Update(&updatedUser).Error
	return
}
