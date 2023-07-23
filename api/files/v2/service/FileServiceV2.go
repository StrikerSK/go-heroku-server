package service

import (
	"github.com/google/uuid"
	fileDomains "go-heroku-server/api/files/v2/domain"
	filePorts "go-heroku-server/api/files/v2/ports"
	"go-heroku-server/api/src/errors"
)

type FileService struct {
	metadataRepository filePorts.IFileMetadataRepositoryV2
	fileRepository     filePorts.IFileRepositoryV2
}

func NewFileService(metadataRepository filePorts.IFileMetadataRepositoryV2, fileRepository filePorts.IFileRepositoryV2) FileService {
	return FileService{
		metadataRepository: metadataRepository,
		fileRepository:     fileRepository,
	}
}

// Function stores files received from the Front-End
func (s FileService) CreateFile(file fileDomains.FileObjectV2) (string, error) {
	id := uuid.New().String()
	file.SetID(id)

	if err := s.metadataRepository.CreateMetadata(file.FileMetadataV2); err != nil {
		return "", err
	}

	if err := s.fileRepository.CreateFile(file.FileEntityV2); err != nil {
		return "", err
	}

	return id, nil
}

// Function provides requested file to the client
func (s FileService) ReadMetadata(fileID, username string) (fileDomains.FileMetadataV2, error) {
	file, err := s.metadataRepository.ReadMetadata(fileID)
	if err != nil {
		return fileDomains.FileMetadataV2{}, err
	} else {
		if file.Username != username {
			return fileDomains.FileMetadataV2{}, errors.NewForbiddenError("forbidden access")
		}
		return file, nil
	}
}

func (s FileService) ReadFiles(username string) ([]fileDomains.FileMetadataV2, error) {
	files, err := s.metadataRepository.ReadAllMetadata(username)
	if err != nil {
		return nil, err
	} else {
		return files, nil
	}
}

// Deletion of file base on userID
func (s FileService) RemoveFile(fileID, username string) error {
	metadata, err := s.metadataRepository.ReadMetadata(fileID)
	if err != nil {
		return err
	}

	if metadata.Username != username {
		return errors.NewForbiddenError("forbidden access")
	}

	err = s.metadataRepository.DeleteMetadata(fileID)
	if err != nil {
		return err
	}

	err = s.fileRepository.DeleteFile(fileID)
	if err != nil {
		return err
	}

	return nil
}

func (s FileService) DownloadFile(fileID, username string) (fileDomains.FileEntityV2, error) {
	metadata, err := s.metadataRepository.ReadMetadata(fileID)
	if err != nil {
		return fileDomains.FileEntityV2{}, err
	}

	if metadata.Username != username {
		return fileDomains.FileEntityV2{}, errors.NewForbiddenError("forbidden access")
	}

	return s.fileRepository.ReadFile(fileID)
}
