package image

type LocationImage struct {
	Id         uint   `json:"id"`
	LocationID uint   `json:"-"`
	FileName   string `json:"filename"`
	FileType   string `json:"fileType"`
	FileData   []byte `json:"-"`
}
