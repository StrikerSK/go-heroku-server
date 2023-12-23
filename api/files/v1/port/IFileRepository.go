package port

import (
	fileDomains "go-heroku-server/api/files/v1/domain"
)

type IFileRepository interface {
	CreateFile(*fileDomains.FileEntity) error
	ReadFile(uint) (fileDomains.FileEntity, error)
	ReadFiles(string) ([]fileDomains.FileEntity, error)
	DeleteFile(uint) error
}
