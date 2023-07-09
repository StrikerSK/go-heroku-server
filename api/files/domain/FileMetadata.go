package fileDomains

import "time"

type FileMetadata struct {
	Id         uint      `json:"id"`
	FileName   string    `json:"fileName"`
	FileType   string    `json:"fileType"`
	FileSize   string    `json:"fileSize"`
	CreateDate time.Time `json:"createDate"`
}
