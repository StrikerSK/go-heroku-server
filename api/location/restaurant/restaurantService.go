package restaurant

import (
	"go-heroku-server/config"

	"encoding/json"
	"log"
	"net/http"
)

type TemporaryName struct {
	Name string "name"
}

func GetRestaurantLocations(w http.ResponseWriter, r *http.Request) {
	var locations []RestaurantLocation
	config.DBConnection.Find(&locations)
	json.NewEncoder(w).Encode(locations)
	log.Println("Retrieved list of restaurant locations")
}

func GetRestaurantByName(w http.ResponseWriter, r *http.Request) {

	var restName TemporaryName
	var restaurant RestaurantLocation

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&restName)
	if err != nil {
		panic(err)
	}
	defer config.DBConnection.Close()
	config.DBConnection.Where("name = ?", restName.Name).First(&restaurant)
	_ = json.NewEncoder(w).Encode(restaurant)
	log.Println("Retrieved restaurant location")
}
