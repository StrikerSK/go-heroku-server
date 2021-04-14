package location

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-heroku-server/api/user"
	"net/http"
)

func EnrichRouteWithLocation(router *mux.Router) {

	subroute := router.PathPrefix("/location").Subrouter()
	subroute.Handle("/save", user.VerifyJwtToken(http.HandlerFunc(controllerAddLocation))).Methods("POST")
	//subroute.HandleFunc("/", addLocation).Methods("GET")
	//subroute.HandleFunc("/", addLocation).Methods("DELETE")

	imageSubroute := subroute.PathPrefix("/image").Subrouter()
	imageSubroute.HandleFunc("/{id}", GetLocationImage).Methods("GET")

	locationsRoute := router.PathPrefix("/location").Subrouter()
	locationsRoute.HandleFunc("/", GetLocations).Methods("GET")

}

func controllerAddLocation(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(user.UserIdContextKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var location Location
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	addLocation(userID, location)
}
