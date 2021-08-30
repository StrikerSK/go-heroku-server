package restaurant

import (
	"go-heroku-server/src/api/location"
)

type RestaurantLocation struct {
	location.Location
	MenuURL string `json:"url"`
}
