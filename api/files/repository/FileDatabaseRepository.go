package fileRepositories

import (
	"fmt"
	fileDomains "go-heroku-server/api/files/domain"
	"go-heroku-server/api/src/errors"
	"gorm.io/gorm"
)

type FileDatabaseRepository struct {
	db *gorm.DB
}

func NewFileDatabaseRepository(db *gorm.DB) FileDatabaseRepository {
	_ = db.AutoMigrate(&fileDomains.FileEntity{})
	return FileDatabaseRepository{
		db: db,
	}
}

func (r FileDatabaseRepository) CreateFile(file *fileDomains.FileEntity) (err error) {
	if err = r.db.Create(&file).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

func (r FileDatabaseRepository) ReadFiles(username string) (files []fileDomains.FileEntity, err error) {
	if err = r.db.Where("username = ?", username).Find(&files).Error; err != nil {
		err = errors.NewDatabaseError(err.Error())
	}
	return
}

func (r FileDatabaseRepository) ReadFile(fileID uint) (file fileDomains.FileEntity, err error) {
	if err = r.db.Where("id = ?", fileID).Find(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fileDomains.FileEntity{}, errors.NewNotFoundError(fmt.Sprintf("file [%d] not found", fileID))
		} else {
			return fileDomains.FileEntity{}, errors.NewDatabaseError(err.Error())
		}
	}
	return
}

func (r FileDatabaseRepository) DeleteFile(file fileDomains.FileEntity) (err error) {
	if err = r.db.Delete(&file).Error; err != nil {
		return errors.NewDatabaseError(err.Error())
	}
	return
}
