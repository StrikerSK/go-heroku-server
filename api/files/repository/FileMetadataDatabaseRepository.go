package fileRepositories

import (
	"fmt"
	fileDomains "go-heroku-server/api/files/domain"
	"go-heroku-server/api/src/errors"
	"gorm.io/gorm"
)

type FileMetadataDatabaseRepository struct {
	db *gorm.DB
}

func NewFileMetadataDatabaseRepository(db *gorm.DB) FileMetadataDatabaseRepository {
	_ = db.AutoMigrate(&fileDomains.FileMetadata{})
	return FileMetadataDatabaseRepository{
		db: db,
	}
}

func (r FileMetadataDatabaseRepository) CreateMetadata(file *fileDomains.FileMetadata) (err error) {
	if err = r.db.Create(&file).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

func (r FileMetadataDatabaseRepository) ReadAll(username string) (files []fileDomains.FileMetadata, err error) {
	if err = r.db.Where("username = ?", username).Find(&files).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

func (r FileMetadataDatabaseRepository) ReadMetadata(fileID uint) (file fileDomains.FileMetadata, err error) {
	if err = r.db.Where("id = ?", fileID).Find(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fileDomains.FileMetadata{}, errors.NewNotFoundError(fmt.Sprintf("file [%d] not found", fileID))
		} else {
			return fileDomains.FileMetadata{}, errors.NewDatabaseError(err.Error())
		}
	}
	return
}

func (r FileMetadataDatabaseRepository) DeleteMetadata(file fileDomains.FileMetadata) (err error) {
	if err = r.db.Delete(&file).Error; err != nil {
		return errors.NewDatabaseError(err.Error())
	}
	return
}
