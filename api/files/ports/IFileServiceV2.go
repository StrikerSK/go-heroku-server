package filePorts

import (
	fileDomains "go-heroku-server/api/files/domain"
)

type IFileServiceV2 interface {
	CreateFile(fileDomains.FileObject) (uint, error)
	ReadFiles(string) ([]fileDomains.FileMetadata, error)
	ReadMetadata(uint, string) (fileDomains.FileMetadata, error)
	DownloadFile(uint, string) (fileDomains.FileEntity, error)
	RemoveFile(uint, string) error
}
