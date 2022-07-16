package locationDomains

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

type UserLocationEntity struct {
	Location
	ImageId  uint   `json:"-"`
	Username string `json:"-"`
}

func (r UserLocationEntity) toString() {
	fmt.Printf("Id: %d, Latitude: %f, Longitutde: %f, Name: %s, Type: %s, Description: %s, ImageId: %d, Username: %s\n",
		r.Id, r.Latitude, r.Longitude, r.Name, r.Type, r.Description, r.ImageId, r.Username)
}
