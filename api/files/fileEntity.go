package files

import (
	"errors"
	"fmt"
	"time"
)

type File struct {
	Id         uint      `json:"id"`
	UserID     uint      `json:"-"`
	FileName   string    `json:"fileName"`
	FileType   string    `json:"fileType"`
	FileData   []byte    `json:"-"`
	FileSize   string    `json:"fileSize"`
	CreateDate time.Time `json:"createDate"`
}

func (f File) toString() {
	fmt.Printf("ID: %d, UserID: %d, FileName: %s, FileType: %s, FileSize: %s \n",
		f.Id, f.UserID, f.FileName, f.FileType, f.FileSize)
}

func (f *File) validateAccess(userID uint) error {
	if f.UserID != userID {
		return errors.New("access denied")
	} else {
		return nil
	}
}
