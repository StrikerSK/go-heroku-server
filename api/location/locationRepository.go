package location

import "go-heroku-server/config"

func createLocation(location Location) {
	config.DBConnection.NewRecord(location)
	config.DBConnection.Create(&location)
}

func readLocation(locationID uint) (location Location, err error) {
	err = config.DBConnection.Where("id = ?", locationID).Find(&location).Error
	return
}

func readAllLocations(userID uint) (locations []Location, err error) {
	err = config.DBConnection.Where("user_id = ?", userID).Find(&locations).Error
	return
}

func deleteLocationFromRepository(location Location) (err error) {
	err = config.DBConnection.Delete(&location).Error
	return
}
