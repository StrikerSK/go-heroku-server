package restaurant

import (
	"go-heroku-server/api/src/responses"
	"net/http"

	"log"
)

func sGetRestaurantByName(name string) responses.IResponse {
	restaurant, err := findByName(name)
	if err != nil {
		log.Printf("Restaurant location [%s] read: %s\n", name, err.Error())
		return responses.NewEmptyResponse(http.StatusNotFound)
	}

	log.Printf("Restaurant location [%s] read: success\n", name)
	return responses.NewResponse(restaurant)
}
