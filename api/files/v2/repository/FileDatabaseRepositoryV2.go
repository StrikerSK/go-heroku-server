package repository

import (
	"fmt"
	fileDomains "go-heroku-server/api/files/v2/domain"
	"go-heroku-server/api/src/errors"
	"go-heroku-server/config"
	"gorm.io/gorm"
)

type FileDatabaseRepository struct {
	db *gorm.DB
}

func NewFileDatabaseRepository() FileDatabaseRepository {
	db := config.GetDatabaseInstance()
	_ = db.AutoMigrate(&fileDomains.FileEntityV2{})
	return FileDatabaseRepository{
		db: db,
	}
}

func (r FileDatabaseRepository) CreateFile(file fileDomains.FileEntityV2) (err error) {
	if err = r.db.Create(&file).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

func (r FileDatabaseRepository) ReadFile(fileID string) (file fileDomains.FileEntityV2, err error) {
	if err = r.db.Where("id = ?", fileID).Find(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fileDomains.FileEntityV2{}, errors.NewNotFoundError(fmt.Sprintf("file [%d] not found", fileID))
		} else {
			return fileDomains.FileEntityV2{}, errors.NewDatabaseError(err.Error())
		}
	}
	return
}

func (r FileDatabaseRepository) DeleteFile(fileID string) (err error) {
	var file fileDomains.FileEntityV2
	if err = r.db.Where("id = ?", fileID).Delete(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError(fmt.Sprintf("file [%d] not found", fileID))
		} else {
			return errors.NewDatabaseError(err.Error())
		}
	}
	return
}
