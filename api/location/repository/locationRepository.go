package locationRepositories

import (
	"fmt"
	"github.com/jinzhu/gorm"
	locationDomains "go-heroku-server/api/location/domain"
	"go-heroku-server/api/src/errors"
)

type LocationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	db.AutoMigrate(&locationDomains.UserLocationEntity{})
	return LocationRepository{
		db: db,
	}
}

func (r LocationRepository) CreateLocation(location *locationDomains.UserLocationEntity) (err error) {
	if err = r.db.Create(&location).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

func (r LocationRepository) ReadLocation(locationID uint) (location locationDomains.UserLocationEntity, err error) {
	if err = r.db.Where("id = ?", locationID).Find(&location).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return locationDomains.UserLocationEntity{}, errors.NewNotFoundError(fmt.Sprintf("location [%d] not found", locationID))
		} else {
			return locationDomains.UserLocationEntity{}, errors.NewDatabaseError(err.Error())
		}
	}
	return
}

func (r LocationRepository) ReadLocations(username string) (locations []locationDomains.UserLocationEntity, err error) {
	if err = r.db.Where("username = ?", username).Find(&locations).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

func (r LocationRepository) UpdateLocation(location locationDomains.UserLocationEntity) (err error) {
	if err = r.db.Save(&location).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.NewNotFoundError(fmt.Sprintf("location [%d] not found", location.Id))
		} else {
			err = errors.NewDatabaseError(err.Error())
		}
	}
	return
}

func (r LocationRepository) DeleteLocation(location locationDomains.UserLocationEntity) (err error) {
	if err = r.db.Delete(&location).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = errors.NewNotFoundError(fmt.Sprintf("location [%d] not found", location.Id))
		} else {
			err = errors.NewDatabaseError(err.Error())
		}
	}
	return
}
