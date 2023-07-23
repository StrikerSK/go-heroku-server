package filePorts

import (
	"go-heroku-server/api/files/v2/domain"
)

type IFileMetadataRepositoryV2 interface {
	CreateMetadata(fileDomains.FileMetadataV2) error
	ReadMetadata(string) (fileDomains.FileMetadataV2, error)
	ReadAllMetadata(string) ([]fileDomains.FileMetadataV2, error)
	DeleteMetadata(string) error
}
