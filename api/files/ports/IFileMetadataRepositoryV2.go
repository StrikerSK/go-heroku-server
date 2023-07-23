package filePorts

import fileDomains "go-heroku-server/api/files/domain"

type IFileMetadataRepositoryV2 interface {
	CreateMetadata(fileDomains.FileMetadataV2) error
	ReadMetadata(uint) (fileDomains.FileMetadataV2, error)
	ReadAllMetadata(string) ([]fileDomains.FileMetadataV2, error)
	DeleteMetadata(uint) error
}
