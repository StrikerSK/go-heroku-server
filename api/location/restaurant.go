package location

import (
	"go-heroku-server/api/types"
	"go-heroku-server/config"

	"encoding/json"
	"log"
	"net/http"
)

type RestaurantLocation struct {
	Id          uint    `json:"id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	MenuURL     string  `json:"url"`
}

func GetRestaurantLocations(w http.ResponseWriter, r *http.Request) {

	var locations []RestaurantLocation

	db, err := config.CreateDatabase()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.Find(&locations)

	json.NewEncoder(w).Encode(locations)
	log.Println("Retrieved list of restaurant locations")

}

func GetReastaurantByName(w http.ResponseWriter, r *http.Request) {

	var restName types.TemporaryName
	var restaurant RestaurantLocation

	db, err := config.CreateDatabase()

	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&restName)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.Where("name = ?", restName.Name).First(&restaurant)

	json.NewEncoder(w).Encode(restaurant)
	log.Println("Retrieved restaurant location")

}
