package restaurant

import (
	"go-heroku-server/src/responses"
	"net/http"

	"log"
)

func sGetRestaurantByName(name string) responses.IResponse {
	restaurant, err := findByName(name)
	if err != nil {
		log.Printf("Restaurant location [%s] read: %v\n", name, err)
		return responses.CreateResponse(http.StatusNotFound, nil)
	}

	log.Printf("Restaurant location [%s] read: success\n", name)
	return responses.CreateResponse(http.StatusOK, restaurant)
}
