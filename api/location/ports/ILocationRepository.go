package locationPorts

import locationDomains "go-heroku-server/api/location/domain"

type ILocationRepository interface {
	CreateLocation(locationDomains.UserLocationEntity) error
	ReadLocation(uint) (locationDomains.UserLocationEntity, error)
	ReadLocations(string) ([]locationDomains.UserLocationEntity, error)
	UpdateLocation(locationDomains.UserLocationEntity) error
	DeleteLocation(locationDomains.UserLocationEntity) error
}
