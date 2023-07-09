package filePorts

import (
	fileDomains "go-heroku-server/api/files/domain"
)

type IFileRepository interface {
	CreateFile(*fileDomains.FileEntity) error
	ReadFile(uint) (fileDomains.FileEntity, error)
	ReadFiles(string) ([]fileDomains.FileEntity, error)
	DeleteFile(fileDomains.FileEntity) error
}
