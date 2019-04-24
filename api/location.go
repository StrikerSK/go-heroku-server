package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"go-heroku-server/config"
)

type Location struct {
	Id          uint    `json:"id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	ImageId     uint    `json:"imageid"`
}

type LocationImage struct {
	Id       uint     `json:"id"`
	Location Location `json:"-"`
	FileName string   `json:"filename"`
	FileType string   `json:"-"`
	FileData []byte   `json:"-"`
}

func GetLocations(w http.ResponseWriter, r *http.Request) {

	var locations []Location

	db, err := config.CreateDatabase()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.Find(&locations)

	json.NewEncoder(w).Encode(locations)
	log.Println("Retrieved list of location")

}

func GetLocationImage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uri, _ := strconv.ParseInt(vars["id"], 10, 64)

	var receivedImage = getImageFromDb(uri)

	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, Content-Length, X-Content-Transfer-Id")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Disposition", "attachment; filename="+receivedImage.FileName)
	w.Header().Set("Content-Type", receivedImage.FileType)

	w.Write(receivedImage.FileData)
}

func getImageFromDb(fileId int64) LocationImage {

	db, err := config.CreateDatabase()
	if err != nil {
		panic(err)
	}

	var image LocationImage

	db.Where("id = ?", fileId).Find(&image)

	return image

}

func AddLocation(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var location Location
	err := decoder.Decode(&location)

	if err != nil {
		panic(err)
	}

	saveLocation(location.Name, location.Description, location.Type, location.Latitude, location.Longitude)
}

func saveLocation(name, description, serviceType string, latitude, longitude float64) {
	db, err := config.CreateDatabase()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	location := Location{Name: name, Description: description, Type: serviceType, Longitude: longitude, Latitude: latitude, ImageId: 1}

	db.NewRecord(location)
	db.Create(&location)

	log.Println("Inserted location with name: " + location.Name + ".")
}
