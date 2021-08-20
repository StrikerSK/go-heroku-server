package image

import "go-heroku-server/api/location"

type LocationImage struct {
	Id       uint                  `json:"id"`
	Location location.UserLocation `json:"-"`
	FileName string                `json:"filename"`
	FileType string                `json:"-"`
	FileData []byte                `json:"-"`
}
