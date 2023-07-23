package filePorts

import (
	fileDomains "go-heroku-server/api/files/domain"
	"go-heroku-server/api/files/v1/domain"
)

type IFileServiceV2 interface {
	CreateFile(fileDomains.FileObjectV2) (uint, error)
	ReadFiles(string) ([]fileDomains.FileMetadataV2, error)
	ReadMetadata(uint, string) (fileDomains.FileMetadataV2, error)
	DownloadFile(uint, string) (domain.FileEntity, error)
	RemoveFile(uint, string) error
}
