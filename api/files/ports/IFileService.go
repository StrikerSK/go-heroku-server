package filePorts

import (
	fileDomains "go-heroku-server/api/files/domain"
)

type IFileService interface {
	CreateFile(fileDomains.FileObject) (uint, error)
	ReadFile(uint, string) (fileDomains.FileObject, error)
	ReadFiles(string) ([]fileDomains.FileObject, error)
	DeleteFile(uint, string) error
}
