package fileServices

import (
	"fmt"
	"github.com/jinzhu/gorm"
	fileDomains "go-heroku-server/api/files/domain"
	filePorts "go-heroku-server/api/files/ports"
	"go-heroku-server/api/src/errors"
	"log"
	"strings"
)

type FileService struct {
	repository filePorts.IFileRepository
}

func NewFileService(repository filePorts.IFileRepository) FileService {
	return FileService{
		repository: repository,
	}
}

//Function stores files received from the Front-End
func (s FileService) CreateFile(fileEntity fileDomains.FileEntity) (uint, error) {
	if err := s.repository.CreateFile(&fileEntity); err != nil {
		return 0, err
	}
	return fileEntity.Id, nil
}

//Function provides requested file to the client
func (s FileService) ReadFile(fileID uint, username string) (fileDomains.FileEntity, error) {
	file, err := s.repository.ReadFile(fileID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fileDomains.FileEntity{}, errors.NewNotFoundError(fmt.Sprintf("file [%d] not found", fileID))
		} else {
			return fileDomains.FileEntity{}, err
		}
	} else {
		if file.Username != username {
			return fileDomains.FileEntity{}, errors.NewForbiddenError("forbidden access")
		}
		return file, nil
	}
}

func (s FileService) ReadFiles(username string) ([]fileDomains.FileEntity, error) {
	files, err := s.repository.ReadFiles(username)
	if err != nil {
		return nil, err
	} else {
		for index := range files {
			fileName := files[index].FileName
			fileName = fileName[:strings.IndexByte(fileName, '.')]
			files[index].FileName = fileName
		}
		log.Printf("File listing: success\n")
		return files, nil
	}
}

//Deletion of file base on userID
func (s FileService) DeleteFile(fileID uint, username string) error {
	_, err := s.ReadFile(fileID, username)
	if err != nil {
		return err
	}

	return s.repository.DeleteFile(fileID)
}