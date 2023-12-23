package domain

import (
	"fmt"
	"time"
)

type FileEntity struct {
	Id         uint      `json:"id"`
	Username   string    `json:"-"`
	FileName   string    `json:"fileName"`
	FileType   string    `json:"fileType"`
	FileData   []byte    `json:"-"`
	FileSize   string    `json:"fileSize"`
	CreateDate time.Time `json:"createDate"`
}

func (r FileEntity) toString() {
	fmt.Printf("ID: %d, Username: %s, FileName: %s, FileType: %s, FileSize: %s \n",
		r.Id, r.Username, r.FileName, r.FileType, r.FileSize)
}
