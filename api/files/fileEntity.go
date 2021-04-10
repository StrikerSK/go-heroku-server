package files

import "time"

type File struct {
	Id         int64     `json:"id"`
	UserID     uint      `json:"-"`
	FileName   string    `json:"fileName"`
	FileType   string    `json:"fileType"`
	FileData   []byte    `json:"-"`
	FileSize   string    `json:"fileSize"`
	CreateDate time.Time `json:"createDate"`
}

type Metadata struct {
	Hello   string
	Testing string
}
