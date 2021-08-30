package restaurant

import (
	"github.com/gorilla/mux"
	"go-heroku-server/config"
	"go-heroku-server/src/responses"
	"net/http"
)

func EnrichRouteWithRestaurants(router *mux.Router) {
	config.InitializeType("Restaurant", &RestaurantLocation{})

	router.Handle("/restaurants", http.HandlerFunc(GetRestaurantLocations)).Methods(http.MethodGet)
	router.Handle("/restaurant", http.HandlerFunc(GetRestaurantByName)).Methods(http.MethodGet)
}

func GetRestaurantLocations(w http.ResponseWriter, r *http.Request) {
	responses.CreateResponse(http.StatusOK, findAll()).WriteResponse(w)
	return
}

func GetRestaurantByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	sGetRestaurantByName(name).WriteResponse(w)
	return
}
