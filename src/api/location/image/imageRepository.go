package image

import "go-heroku-server/config"

func readImage(imageID int64) (image LocationImage, err error) {
	err = config.GetDatabaseInstance().Where("id = ?", imageID).Find(&image).Error
	return
}
