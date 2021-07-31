package location

import "fmt"

type Location struct {
	Id          uint    `json:"-"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	ImageId     uint    `json:"imageId"`
	UserID      uint    `json:"-"`
}

func (r Location) toString() {
	fmt.Printf("Id: %d, Latitude: %f, Longitutde: %f, Name: %s, Type: %s, Description: %s, ImageId: %d, UserID: %d\n",
		r.Id, r.Latitude, r.Longitude, r.Name, r.Type, r.Description, r.ImageId, r.UserID)
}

type LocationImage struct {
	Id       uint     `json:"id"`
	Location Location `json:"-"`
	FileName string   `json:"filename"`
	FileType string   `json:"-"`
	FileData []byte   `json:"-"`
}
