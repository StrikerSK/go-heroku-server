package repository

import (
	"fmt"
	fileDomains "go-heroku-server/api/files/v2/domain"
	"go-heroku-server/api/src/errors"
	"go-heroku-server/config"
	"gorm.io/gorm"
)

type FileMetadataRepository struct {
	db *gorm.DB
}

func NewFileMetadataRepository() FileMetadataRepository {
	db := config.GetDatabaseInstance()
	_ = db.AutoMigrate(&fileDomains.FileMetadataV2{})
	return FileMetadataRepository{
		db: db,
	}
}

func (r FileMetadataRepository) CreateMetadata(file fileDomains.FileMetadataV2) (err error) {
	if err = r.db.Create(&file).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

func (r FileMetadataRepository) ReadMetadata(fileID string) (metadata fileDomains.FileMetadataV2, err error) {
	if err = r.db.Where("id = ?", fileID).Find(&metadata).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fileDomains.FileMetadataV2{}, errors.NewNotFoundError(fmt.Sprintf("file [%d] not found", fileID))
		} else {
			return fileDomains.FileMetadataV2{}, errors.NewDatabaseError(err.Error())
		}
	}
	return
}

func (r FileMetadataRepository) ReadAllMetadata(username string) (metadata []fileDomains.FileMetadataV2, err error) {
	if err = r.db.Where("username = ?", username).Find(&metadata).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

func (r FileMetadataRepository) DeleteMetadata(fileID string) (err error) {
	var file fileDomains.FileMetadataV2
	if err = r.db.Where("id = ?", fileID).Delete(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError(fmt.Sprintf("Metadat [%s] not found", fileID))
		} else {
			return errors.NewDatabaseError(err.Error())
		}
	}
	return
}
