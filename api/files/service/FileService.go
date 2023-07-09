package fileServices

import (
	fileDomains "go-heroku-server/api/files/domain"
	filePorts "go-heroku-server/api/files/ports"
	"go-heroku-server/api/src/errors"
)

type FileService struct {
	repository         filePorts.IFileRepository
	metadataRepository filePorts.IFileMetadataRepository
}

func NewFileService(repository filePorts.IFileRepository, metadataRepository filePorts.IFileMetadataRepository) FileService {
	return FileService{
		repository:         repository,
		metadataRepository: metadataRepository,
	}
}

// Function stores files received from the Front-End
func (s FileService) CreateFile(fileContent fileDomains.FileEntity, fileMetadata fileDomains.FileMetadata) (uint, error) {
	if err := s.metadataRepository.CreateMetadata(&fileMetadata); err != nil {
		return 0, err
	}

	if err := s.repository.CreateFile(&fileContent); err != nil {
		return 0, err
	}

	return fileContent.Id, nil
}

// Function provides requested file to the client
func (s FileService) ReadFile(fileID uint, username string) (fileDomains.FileEntity, error) {
	file, err := s.repository.ReadFile(fileID)
	if err != nil {
		return fileDomains.FileEntity{}, err
	} else {
		if file.Username != username {
			return fileDomains.FileEntity{}, errors.NewForbiddenError("forbidden access")
		}
		return file, nil
	}
}

func (s FileService) ReadFiles(username string) ([]fileDomains.FileMetadata, error) {
	files, err := s.metadataRepository.ReadAll(username)
	if err != nil {
		return nil, err
	} else {
		return files, nil
	}
}

// Deletion of file base on userID
func (s FileService) DeleteFile(fileID uint, username string) error {
	file, err := s.ReadFile(fileID, username)
	if err != nil {
		return err
	}

	return s.repository.DeleteFile(file)
}
