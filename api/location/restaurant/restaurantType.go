package restaurant

import "go-heroku-server/api/location"

type RestaurantLocation struct {
	location.Location
	MenuURL string `json:"url"`
}
