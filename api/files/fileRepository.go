package files

import (
	"go-heroku-server/config"
)

func createFile(file File) {
	instance := config.GetDatabaseInstance()
	instance.NewRecord(file)
	instance.Create(&file)
}

func getAll(userID uint) (files []File) {
	config.GetDatabaseInstance().Where("user_id = ?", userID).Find(&files)
	return
}

func getFile(fileId uint) (file File, err error) {
	err = config.GetDatabaseInstance().Where("id = ?", fileId).Find(&file).Error
	return
}

func deleteFile(fileID uint) (err error) {
	var file File
	err = config.GetDatabaseInstance().Where("id = ?", fileID).Delete(&file).Error
	return
}
