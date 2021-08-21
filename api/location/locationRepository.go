package location

import (
	"go-heroku-server/config"
)

func createLocation(location UserLocation) {
	config.GetDatabaseInstance().NewRecord(location)
	config.GetDatabaseInstance().Create(&location)
}

func readLocation(locationID uint) (location UserLocation, err error) {
	err = config.GetDatabaseInstance().Where("id = ?", locationID).Find(&location).Error
	return
}

func readAllLocations(userID uint) (locations []UserLocation, err error) {
	err = config.GetDatabaseInstance().Where("user_id = ?", userID).Find(&locations).Error
	return
}

func updateLocationInRepository(location UserLocation) (err error) {
	err = config.GetDatabaseInstance().Save(&location).Error
	return
}

func deleteLocationFromRepository(location UserLocation) (err error) {
	err = config.GetDatabaseInstance().Delete(&location).Error
	return
}
