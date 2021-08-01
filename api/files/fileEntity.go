package files

import (
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

func (r File) toString() {
	fmt.Printf("ID: %d, UserID: %d, FileName: %s, FileType: %s, FileSize: %s \n",
		r.Id, r.UserID, r.FileName, r.FileType, r.FileSize)
}
