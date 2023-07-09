package filePorts

import (
	fileDomains "go-heroku-server/api/files/domain"
)

type IFileMetadataRepository interface {
	CreateMetadata(*fileDomains.FileMetadata) error
	ReadMetadata(uint) (fileDomains.FileMetadata, error)
	ReadAll(string) ([]fileDomains.FileMetadata, error)
	DeleteMetadata(fileDomains.FileMetadata) error
}
