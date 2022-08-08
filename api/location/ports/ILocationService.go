package locationPorts

import (
	locationDomains "go-heroku-server/api/location/domain"
)

type ILocationService interface {
	CreateLocation(locationDomains.UserLocationEntity) (uint, error)
	ReadLocation(uint, string) (locationDomains.UserLocationEntity, error)
	ReadLocations(string) ([]locationDomains.UserLocationEntity, error)
	UpdateLocation(locationDomains.UserLocationEntity) error
	DeleteLocation(uint, string) error
}
