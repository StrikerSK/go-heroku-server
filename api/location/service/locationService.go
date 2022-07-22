package locationServcices

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"go-heroku-server/api/location/domain"
	locationPorts "go-heroku-server/api/location/ports"
	errors2 "go-heroku-server/api/src/errors"
	"log"
)

type LocationService struct {
	repository locationPorts.ILocationRepository
}

func NewLocationService(repository locationPorts.ILocationRepository) LocationService {
	return LocationService{
		repository: repository,
	}
}

func (s LocationService) CreateLocation(location locationDomains.UserLocationEntity) (uint, error) {
	return s.repository.CreateLocation(location)
}

func (s LocationService) DeleteLocation(locationID uint, username string) error {
	persistedLocation, err := s.ReadLocation(locationID, username)
	if err != nil {
		log.Printf("UserLocationEntity [%d] delete: %s", locationID, err.Error())
		return err
	}

	err = s.repository.DeleteLocation(persistedLocation)
	if err != nil {
		log.Printf("UserLocationEntity [%d] delete: %s", locationID, err.Error())
		return err
	}

	return nil
}

func (s LocationService) UpdateLocation(updatedLocation locationDomains.UserLocationEntity) error {
	_, err := s.ReadLocation(updatedLocation.Id, updatedLocation.Username)
	if err != nil {
		return err
	}

	if err = s.repository.UpdateLocation(updatedLocation); err != nil {
		log.Printf("Location [%d] edit: %s", updatedLocation.Id, err.Error())
		return err
	} else {
		//log.Printf("Location [%d] edit: success", locationID)
		return nil
	}
}

func (s LocationService) ReadLocation(locationID uint, username string) (locationDomains.UserLocationEntity, error) {
	var location, err = s.repository.ReadLocation(locationID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("Location [%d] read: %v", locationID, err.Error())
			return locationDomains.UserLocationEntity{}, errors2.NewNotFoundError(fmt.Sprintf("location [%d] not found", locationID))
		} else {
			log.Printf("Location [%d] read: %v", locationID, err)
			return locationDomains.UserLocationEntity{}, err
		}
	} else {
		if location.Username != username {
			log.Printf("Location [%d] read: access denied", locationID)
			return locationDomains.UserLocationEntity{}, errors2.NewForbiddenError("access forbidden")
		} else {
			return location, nil
		}
	}
}

func (s LocationService) ReadLocations(username string) ([]locationDomains.UserLocationEntity, error) {
	return s.repository.ReadLocations(username)
}
