package filePorts

import (
	fileDomains "go-heroku-server/api/files/domain"
)

type IFileRepository interface {
	CreateFile(fileDomains.FileEntity) error
	GetAll(string) ([]fileDomains.FileEntity, error)
	GetFile(uint) (fileDomains.FileEntity, error)
	DeleteFile(uint) error
}
