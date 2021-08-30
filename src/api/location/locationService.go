package location

import (
	"go-heroku-server/src/responses"
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
		log.Printf("UserLocation [%d] delete: %v\n", locationID, err)
		return responses.CreateResponse(http.StatusBadRequest, nil)
	}

	if err = persistedLocation.validateAccess(userID); err != nil {
		log.Printf("UserLocation [%d] delete: %v\n", locationID, err)
		return responses.CreateResponse(http.StatusForbidden, nil)
	}

	if err = deleteLocationFromRepository(persistedLocation); err != nil {
		log.Printf("UserLocation [%d] delete: %v\n", locationID, err)
		return responses.CreateResponse(http.StatusBadRequest, nil)
	}

	return nil
}

func editLocation(userID, locationID uint, updatedLocation UserLocation) responses.IResponse {
	persistedLocation, err := readLocation(locationID)
	if err != nil {
		log.Printf("UserLocation [%d] edit: %v\n", locationID, err)
		return responses.CreateResponse(http.StatusNotFound, nil)
	}

	if err = persistedLocation.validateAccess(userID); err != nil {
		log.Printf("UserLocation [%d] edit: %v\n", locationID, err)
		return responses.CreateResponse(http.StatusForbidden, nil)
	}

	updatedLocation.Id = locationID
	updatedLocation.UserID = persistedLocation.UserID

	if err = updateLocationInRepository(updatedLocation); err != nil {
		log.Printf("Location [%d] edit: %v\n", locationID, err)
		return responses.CreateResponse(http.StatusBadRequest, nil)
	}

	log.Printf("Location [%d] edit: success\n", locationID)
	return responses.CreateResponse(http.StatusOK, nil)
}

func retrieveLocation(userID, locationID uint) (res responses.IResponse) {
	var persistedLocation, err = readLocation(locationID)
	if err != nil {
		log.Printf("Location [%d] read: %v\n", locationID, err)
		res = responses.CreateResponse(http.StatusNotFound, nil)
		return
	}

	if err = persistedLocation.validateAccess(userID); err != nil {
		log.Printf("UserLocation [%d] delete: %v\n", locationID, err)
		return responses.CreateResponse(http.StatusForbidden, nil)
	}

	res = responses.CreateResponse(http.StatusOK, persistedLocation)
	return
}

func getAllLocations(userID uint) (res responses.IResponse) {
	if locations, err := readAllLocations(userID); err != nil {
		log.Printf("Locations reading: %v\n", err)
		res = responses.CreateResponse(http.StatusBadRequest, nil)
	} else {
		res = responses.CreateResponse(http.StatusOK, locations)
	}
	return
}
