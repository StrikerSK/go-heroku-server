package filePorts

import (
	fileDomains "go-heroku-server/api/files/domain"
)

type IFileService interface {
	CreateFile(fileDomains.FileEntity) (uint, error)
	ReadFile(uint, string) (fileDomains.FileEntity, error)
	ReadFiles(string) ([]fileDomains.FileEntity, error)
	DeleteFile(uint, string) error
}
