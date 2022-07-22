package locationRepositories

import (
	"github.com/jinzhu/gorm"
	locationDomains "go-heroku-server/api/location/domain"
)

type LocationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	db.AutoMigrate(&locationDomains.Location{})
	return LocationRepository{
		db: db,
	}
}

func (r LocationRepository) CreateLocation(location locationDomains.UserLocationEntity) (err error) {
	err = r.db.Create(&location).Error
	return
}

func (r LocationRepository) ReadLocation(locationID uint) (location locationDomains.UserLocationEntity, err error) {
	err = r.db.Where("id = ?", locationID).Find(&location).Error
	return
}

func (r LocationRepository) ReadLocations(username string) (locations []locationDomains.UserLocationEntity, err error) {
	err = r.db.Where("username = ?", username).Find(&locations).Error
	return
}

func (r LocationRepository) UpdateLocation(location locationDomains.UserLocationEntity) (err error) {
	err = r.db.Save(&location).Error
	return
}

func (r LocationRepository) DeleteLocation(location locationDomains.UserLocationEntity) (err error) {
	err = r.db.Delete(&location).Error
	return
}
