package location

import (
	"errors"
	"fmt"
)

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

func (ul UserLocation) toString() {
	fmt.Printf("Id: %d, Latitude: %f, Longitutde: %f, Name: %s, Type: %s, Description: %s, ImageId: %d, UserID: %d\n",
		ul.Id, ul.Latitude, ul.Longitude, ul.Name, ul.Type, ul.Description, ul.ImageId, ul.UserID)
}

func (ul *UserLocation) validateAccess(userID uint) error {
	if ul.UserID != userID {
		return errors.New("access denied")
	} else {
		return nil
	}
}
