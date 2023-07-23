package filePorts

import (
	fileDomains "go-heroku-server/api/files/v2/domain"
)

type IFileServiceV2 interface {
	CreateFile(fileDomains.FileObjectV2) (string, error)
	ReadFiles(string) ([]fileDomains.FileMetadataV2, error)
	ReadMetadata(string, string) (fileDomains.FileMetadataV2, error)
	DownloadFile(string, string) (fileDomains.FileEntityV2, error)
	RemoveFile(string, string) error
}
