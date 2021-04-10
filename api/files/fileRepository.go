package files

import (
	"go-heroku-server/config"
)

func createFile(file File) {
	config.DBConnection.NewRecord(file)
	config.DBConnection.Create(&file)
}

func getAll(userID interface{}) (files []File) {
	config.DBConnection.Where("user_id = ?", userID).Find(&files)
	return
}

func getFile(fileId int64) (file File, err error) {
	err = config.DBConnection.Where("id = ?", fileId).Find(&file).Error
	return
}

func deleteFile(fileID interface{}) (file File, err error) {
	err = config.DBConnection.Where("user_id = ?", fileID).Delete(&file).Error
	return
}
