package location

import "fmt"

type Location struct {
	Id          uint    `json:"id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	URL         string  `json:"url,omitempty"`
}

type UserLocation struct {
	Location
	ImageId uint `json:"-"`
	UserID  uint `json:"-"`
}

func (r UserLocation) toString() {
	fmt.Printf("Id: %d, Latitude: %f, Longitutde: %f, Name: %s, Type: %s, Description: %s, ImageId: %d, UserID: %d\n",
		r.Id, r.Latitude, r.Longitude, r.Name, r.Type, r.Description, r.ImageId, r.UserID)
}
