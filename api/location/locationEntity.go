package location

type Location struct {
	Id          uint    `json:"id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	ImageId     uint    `json:"imageId"`
	UserID      uint    `json:"-"`
}

type LocationImage struct {
	Id       uint     `json:"id"`
	Location Location `json:"-"`
	FileName string   `json:"filename"`
	FileType string   `json:"-"`
	FileData []byte   `json:"-"`
}
