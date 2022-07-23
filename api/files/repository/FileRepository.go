package fileRepositories

import (
	"github.com/jinzhu/gorm"
	fileDomains "go-heroku-server/api/files/domain"
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

func (r FileRepository) CreateFile(file *fileDomains.FileEntity) error {
	return r.db.Create(&file).Error
}

func (r FileRepository) ReadFiles(username string) (files []fileDomains.FileEntity, err error) {
	err = r.db.Where("username = ?", username).Find(&files).Error
	return
}

func (r FileRepository) ReadFile(fileId uint) (file fileDomains.FileEntity, err error) {
	err = r.db.Where("id = ?", fileId).Find(&file).Error
	return
}

func (r FileRepository) DeleteFile(fileID uint) (err error) {
	var file fileDomains.FileEntity
	err = r.db.Where("id = ?", fileID).Delete(&file).Error
	return
}
