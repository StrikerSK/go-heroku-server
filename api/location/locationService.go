package location

import (
	"errors"
	"github.com/gorilla/mux"
	"go-heroku-server/api/src"
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

func deleteLocation(userID, locationID uint) src.IResponse {
	persistedLocation, err := readLocation(locationID)
	if err != nil {
		return src.NewErrorResponse(http.StatusBadRequest, err)
	}

	if persistedLocation.UserID != userID {
		return src.NewErrorResponse(http.StatusForbidden, errors.New("user cannot access requested location"))
	}

	if err = deleteLocationFromRepository(persistedLocation); err != nil {
		return src.NewErrorResponse(http.StatusBadRequest, err)
	}

	return nil
}

func editLocation(userID, locationID uint, updatedLocation Location) src.IResponse {
	persistedLocation, err := readLocation(locationID)
	if err != nil {
		return src.NewErrorResponse(http.StatusNotFound, err)
	}

	if persistedLocation.UserID != userID {
		customError := errors.New("user were accessing unowned todo")
		return src.NewErrorResponse(http.StatusForbidden, customError)
	}

	updatedLocation.Id = locationID
	updatedLocation.UserID = persistedLocation.UserID

	if err = updateLocationInRepository(updatedLocation); err != nil {
		return src.NewErrorResponse(http.StatusBadRequest, err)
	}

	return src.NewEmptyResponse(http.StatusOK)
}

func retrieveLocation(userID, locationID uint) (res src.IResponse) {
	var location, err = readLocation(locationID)
	if err != nil {
		res = src.NewErrorResponse(http.StatusNotFound, err)
		return
	}

	if location.UserID != userID {
		res = src.NewErrorResponse(http.StatusForbidden, err)
		return
	}

	res = src.NewResponse(location)
	return
}

func getAllLocations(userID uint) (res src.IResponse) {
	if locations, err := readAllLocations(userID); err != nil {
		res = src.NewErrorResponse(http.StatusBadRequest, err)
	} else {
		res = src.NewResponse(locations)
	}
	return
}
