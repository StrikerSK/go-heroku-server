package filePorts

import fileDomains "go-heroku-server/api/files/domain"

type IFileMetadataRepositoryV2 interface {
	CreateMetadata(fileDomains.FileMetadata) error
	ReadMetadata(uint) (fileDomains.FileMetadata, error)
	ReadAllMetadata(string) ([]fileDomains.FileMetadata, error)
	DeleteMetadata(uint) error
}
