package filePorts

import (
	"go-heroku-server/api/files/v1/domain"
	fileDomains2 "go-heroku-server/api/files/v2/domain"
)

type IFileServiceV2 interface {
	CreateFile(fileDomains2.FileObjectV2) (uint, error)
	ReadFiles(string) ([]fileDomains2.FileMetadataV2, error)
	ReadMetadata(uint, string) (fileDomains2.FileMetadataV2, error)
	DownloadFile(uint, string) (domain.FileEntity, error)
	RemoveFile(uint, string) error
}
