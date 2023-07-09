package fileRepositories

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

type FileSystemRepository struct {
	directory string
}

func NewFileSystemRepository() FileSystemRepository {
	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Failed to get user's home directory:", err)
		os.Exit(1)
	}
	homeDir := usr.HomeDir

	// Create the directory if it doesn't exist
	directory := filepath.Join(homeDir, "attachments")
	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		fmt.Println("Failed to create directory:", err)
		os.Exit(1)
	}

	return FileSystemRepository{
		directory: directory,
	}
}
