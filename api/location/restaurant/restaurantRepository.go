package restaurant

import (
	"go-heroku-server/config"
)

func findAll() (restaurants RestaurantLocation) {
	config.GetDatabaseInstance().Find(&restaurants)
	return
}

func findByName(restaurantName string) (restaurant RestaurantLocation, err error) {
	err = config.GetDatabaseInstance().Where("name = ?", restaurantName).First(&restaurant).Error
	return
}
