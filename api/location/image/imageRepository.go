package image

import "go-heroku-server/config"

func getImageFromDb(fileId int64) LocationImage {
	var image LocationImage
	config.DBConnection.Where("id = ?", fileId).Find(&image)
	return image
}
