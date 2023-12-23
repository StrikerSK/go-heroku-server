package repository

import (
	"fmt"
	fileDomains "go-heroku-server/api/files/v2/domain"
	"go-heroku-server/api/src/errors"
	"os"
	"os/user"
	"path/filepath"
)

type FileSystemRepositoryV2 struct {
	directory string
}

func NewFileSystemRepository() FileSystemRepositoryV2 {
	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Failed to get user's home directory:", err)
		os.Exit(1)
	}
	homeDir := usr.HomeDir

	// Create the directory if it doesn't exist
	directory := filepath.Join(homeDir, "custom_attachments")
	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to create directory:", err)
		os.Exit(1)
	}

	return FileSystemRepositoryV2{
		directory: directory,
	}
}

func (fs FileSystemRepositoryV2) CreateFile(file fileDomains.FileEntityV2) error {
	filePath := fs.directory + "/" + file.Id

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return errors.NewNotFoundError("file already exists")
	}

	// Create the file
	err := os.WriteFile(filePath, []byte(file.FileData), 0644)
	if err != nil {
		return errors.NewNotFoundError(fmt.Sprintf("failed to create file: %v", err))
	}

	return nil
}

func (fs FileSystemRepositoryV2) ReadFile(filename string) (fileDomains.FileEntityV2, error) {
	filePath := fs.directory + "/" + filename

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fileDomains.FileEntityV2{}, errors.NewNotFoundError(fmt.Sprintf("failed to read file: %v", err))
	}

	return fileDomains.FileEntityV2{
		Id:       filename,
		FileData: content,
	}, nil
}

func (fs FileSystemRepositoryV2) DeleteFile(filename string) error {
	filePath := fs.directory + "/" + filename

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errors.NewNotFoundError("file not found")
	}

	// Delete the file
	err := os.Remove(filePath)
	if err != nil {
		return errors.NewNotFoundError(fmt.Sprintf("failed to delete file: %v", err))
	}

	return nil
}
