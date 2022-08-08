package locationServcices

import (
	"go-heroku-server/api/location/domain"
	locationPorts "go-heroku-server/api/location/ports"
	"go-heroku-server/api/src/errors"
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
	if err := s.repository.CreateLocation(&location); err != nil {
		return 0, err
	} else {
		return location.Id, nil
	}
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
	location, err := s.repository.ReadLocation(locationID)
	if err != nil {
		return location, err
	} else {
		if location.Username != username {
			log.Printf("Location [%d] read: access denied", locationID)
			return locationDomains.UserLocationEntity{}, errors.NewForbiddenError("access denied")
		} else {
			return location, nil
		}
	}
}

func (s LocationService) ReadLocations(username string) ([]locationDomains.UserLocationEntity, error) {
	return s.repository.ReadLocations(username)
}
