package locationPorts

import (
	locationDomains "go-heroku-server/api/location/domain"
)

type ILocationService interface {
	CreateLocation(locationDomains.UserLocationEntity) error
	ReadLocations(string) ([]locationDomains.UserLocationEntity, error)
	GetLocation(uint, string) (locationDomains.UserLocationEntity, error)
	UpdateLocation(locationDomains.UserLocationEntity) error
	DeleteLocation(locationID uint, userID string) error
}
