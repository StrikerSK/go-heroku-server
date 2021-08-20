package location

import (
	"go-heroku-server/api/src/responses"
	"log"
	"net/http"
)

func addLocation(userID uint, location UserLocation) {
	location.UserID = userID
	createLocation(location)
}

func deleteLocation(userID, locationID uint) responses.IResponse {
	persistedLocation, err := readLocation(locationID)
	if err != nil {
		log.Printf("UserLocation [%d] delete: %s", locationID, err.Error())
		return responses.NewErrorResponse(http.StatusBadRequest, err)
	}

	if persistedLocation.UserID != userID {
		log.Printf("UserLocation [%d] delete: access denied", locationID)
		return responses.NewEmptyResponse(http.StatusForbidden)
	}

	if err = deleteLocationFromRepository(persistedLocation); err != nil {
		log.Printf("UserLocation [%d] delete: %s", locationID, err.Error())
		return responses.NewErrorResponse(http.StatusBadRequest, err)
	}

	return nil
}

func editLocation(userID, locationID uint, updatedLocation UserLocation) responses.IResponse {
	persistedLocation, err := readLocation(locationID)
	if err != nil {
		log.Printf("UserLocation [%d] edit: %s", locationID, err.Error())
		return responses.NewErrorResponse(http.StatusNotFound, err)
	}

	if persistedLocation.UserID != userID {
		log.Printf("UserLocation [%d] edit: access denied", locationID)
		return responses.NewEmptyResponse(http.StatusForbidden)
	}

	updatedLocation.Id = locationID
	updatedLocation.UserID = persistedLocation.UserID

	if err = updateLocationInRepository(updatedLocation); err != nil {
		log.Printf("UserLocation [%d] edit: %s", locationID, err.Error())
		return responses.NewErrorResponse(http.StatusBadRequest, err)
	}

	return responses.NewEmptyResponse(http.StatusOK)
}

func retrieveLocation(userID, locationID uint) (res responses.IResponse) {
	var location, err = readLocation(locationID)
	if err != nil {
		log.Printf("UserLocation [%d] read: %s", locationID, err.Error())
		res = responses.NewEmptyResponse(http.StatusNotFound)
		return
	}

	if location.UserID != userID {
		log.Printf("UserLocation [%d] read: access denied", locationID)
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
