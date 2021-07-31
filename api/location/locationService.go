package location

import (
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

func retrieveLocation(userID, locationID uint) (*Location, *src.RequestError) {
	var location, err = readLocation(locationID)
	if err != nil {
		log.Printf(err.Error() + " for id: " + strconv.Itoa(int(locationID)))
		return nil, &src.RequestError{
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
	}

	if location.UserID != userID {
		log.Printf("access denied for file id: " + strconv.Itoa(int(locationID)))
		return nil, &src.RequestError{
			StatusCode: http.StatusForbidden,
			Err:        err,
		}
	}

	return &location, nil
}
