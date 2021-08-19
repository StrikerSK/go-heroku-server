package location

import (
	"github.com/gorilla/mux"
	"go-heroku-server/api/src/responses"
	"log"
	"net/http"
	"strconv"

	"go-heroku-server/config"
)

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
	var image LocationImage
	config.DBConnection.Where("id = ?", fileId).Find(&image)
	return image
}

func addLocation(userID uint, location Location) {
	location.UserID = userID
	createLocation(location)
}

func deleteLocation(userID, locationID uint) responses.IResponse {
	persistedLocation, err := readLocation(locationID)
	if err != nil {
		log.Printf("Location [%d] delete: %s", locationID, err.Error())
		return responses.NewErrorResponse(http.StatusBadRequest, err)
	}

	if persistedLocation.UserID != userID {
		log.Printf("Location [%d] delete: access denied", locationID)
		return responses.NewEmptyResponse(http.StatusForbidden)
	}

	if err = deleteLocationFromRepository(persistedLocation); err != nil {
		log.Printf("Location [%d] delete: %s", locationID, err.Error())
		return responses.NewErrorResponse(http.StatusBadRequest, err)
	}

	return nil
}

func editLocation(userID, locationID uint, updatedLocation Location) responses.IResponse {
	persistedLocation, err := readLocation(locationID)
	if err != nil {
		log.Printf("Location [%d] edit: %s", locationID, err.Error())
		return responses.NewErrorResponse(http.StatusNotFound, err)
	}

	if persistedLocation.UserID != userID {
		log.Printf("Location [%d] edit: access denied", locationID)
		return responses.NewEmptyResponse(http.StatusForbidden)
	}

	updatedLocation.Id = locationID
	updatedLocation.UserID = persistedLocation.UserID

	if err = updateLocationInRepository(updatedLocation); err != nil {
		log.Printf("Location [%d] edit: %s", locationID, err.Error())
		return responses.NewErrorResponse(http.StatusBadRequest, err)
	}

	return responses.NewEmptyResponse(http.StatusOK)
}

func retrieveLocation(userID, locationID uint) (res responses.IResponse) {
	var location, err = readLocation(locationID)
	if err != nil {
		log.Printf("Location [%d] read: %s", locationID, err.Error())
		res = responses.NewEmptyResponse(http.StatusNotFound)
		return
	}

	if location.UserID != userID {
		log.Printf("Location [%d] read: access denied", locationID)
		res = responses.NewEmptyResponse(http.StatusForbidden)
		return
	}

	res = responses.NewResponse(location)
	return
}

func getAllLocations(userID uint) (res responses.IResponse) {
	if locations, err := readAllLocations(userID); err != nil {
		log.Printf("Locations reading: %s", err.Error())
		res = responses.NewErrorResponse(http.StatusBadRequest, err)
	} else {
		res = responses.NewResponse(locations)
	}
	return
}
