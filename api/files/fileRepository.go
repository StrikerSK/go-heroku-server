package files

import (
	"go-heroku-server/config"
)

func createFile(file File) {
	config.DBConnection.NewRecord(file)
	config.DBConnection.Create(&file)
}

func getAll(userID uint) (files []File) {
	config.DBConnection.Where("user_id = ?", userID).Find(&files)
	return
}

func getFile(fileId uint) (file File, err error) {
	err = config.DBConnection.Where("id = ?", fileId).Find(&file).Error
	return
}

func deleteFile(fileID uint) (err error) {
	var file File
	err = config.DBConnection.Where("id = ?", fileID).Delete(&file).Error
	return
}
