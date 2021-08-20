package location

import "go-heroku-server/config"

func createLocation(location UserLocation) {
	config.DBConnection.NewRecord(location)
	config.DBConnection.Create(&location)
}

func readLocation(locationID uint) (location UserLocation, err error) {
	err = config.DBConnection.Where("id = ?", locationID).Find(&location).Error
	return
}

func readAllLocations(userID uint) (locations []UserLocation, err error) {
	err = config.DBConnection.Where("user_id = ?", userID).Find(&locations).Error
	return
}

func updateLocationInRepository(location UserLocation) (err error) {
	err = config.DBConnection.Save(&location).Error
	return
}

func deleteLocationFromRepository(location UserLocation) (err error) {
	err = config.DBConnection.Delete(&location).Error
	return
}
