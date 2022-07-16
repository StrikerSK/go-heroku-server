package restaurant

import (
	"go-heroku-server/api/location/domain"
)

type RestaurantLocation struct {
	locationDomains.Location
	MenuURL string `json:"url"`
}
