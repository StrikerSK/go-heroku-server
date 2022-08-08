package fileRepositories

import (
	"fmt"
	"github.com/jinzhu/gorm"
	fileDomains "go-heroku-server/api/files/domain"
	"go-heroku-server/api/src/errors"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	db.AutoMigrate(&fileDomains.FileEntity{})
	return FileRepository{
		db: db,
	}
}

func (r FileRepository) CreateFile(file *fileDomains.FileEntity) (err error) {
	if err = r.db.Create(&file).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

func (r FileRepository) ReadFiles(username string) (files []fileDomains.FileEntity, err error) {
	if err = r.db.Where("username = ?", username).Find(&files).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

func (r FileRepository) ReadFile(fileID uint) (file fileDomains.FileEntity, err error) {
	if err = r.db.Where("id = ?", fileID).Find(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fileDomains.FileEntity{}, errors.NewNotFoundError(fmt.Sprintf("file [%d] not found", fileID))
		} else {
			return fileDomains.FileEntity{}, errors.NewDatabaseError(err.Error())
		}
	}
	return
}

func (r FileRepository) DeleteFile(fileID uint) (err error) {
	var file fileDomains.FileEntity
	if err = r.db.Where("id = ?", fileID).Delete(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError(fmt.Sprintf("file [%d] not found", fileID))
		} else {
			return errors.NewDatabaseError(err.Error())
		}
	}
	return
}
