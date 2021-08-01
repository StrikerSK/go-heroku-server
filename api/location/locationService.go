package location

import (
	"errors"
	"github.com/gorilla/mux"
	"go-heroku-server/api/src"
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

func deleteLocation(userID, locationID uint) src.IResponse {
	persistedLocation, err := readLocation(locationID)
	if err != nil {
		return src.RequestError{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	if persistedLocation.UserID != userID {
		return src.RequestError{
			StatusCode: http.StatusForbidden,
			Err:        errors.New("user cannot access requested location"),
		}
	}

	if err = deleteLocationFromRepository(persistedLocation); err != nil {
		return src.RequestError{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	return nil
}

func retrieveLocation(userID, locationID uint) (res src.IResponse) {
	var location, err = readLocation(locationID)
	if err != nil {
		res = src.RequestError{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
		return
	}

	if location.UserID != userID {
		log.Printf("access denied for file id: " + strconv.Itoa(int(locationID)))
		res = src.RequestError{
			StatusCode: http.StatusForbidden,
			Err:        err,
		}
		return
	}

	res = src.ResponseImpl{
		Data: location,
	}
	return
}

func getAllLocations(userID uint) (res src.IResponse) {
	if locations, err := readAllLocations(userID); err != nil {
		res = src.RequestError{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	} else {
		res = src.ResponseImpl{
			Data: locations,
		}
	}
	return
}
