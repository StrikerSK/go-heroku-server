package location

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-heroku-server/api/user"
	"go-heroku-server/config"
	"net/http"
	"strconv"
)

const locationContextKey = "locationID"

func EnrichRouteWithLocation(router *mux.Router) {

	subroute := router.PathPrefix("/location").Subrouter()
	subroute.Handle("/save", user.VerifyJwtToken(http.HandlerFunc(controllerAddLocation))).Methods("POST")
	subroute.Handle("/{id}", user.VerifyJwtToken(ResolveLocationID(http.HandlerFunc(controllerGetLocation)))).Methods("GET")
	subroute.Handle("/{id}", user.VerifyJwtToken(ResolveLocationID(http.HandlerFunc(controllerDeleteLocation)))).Methods("DELETE")

	imageSubroute := subroute.PathPrefix("/image").Subrouter()
	imageSubroute.HandleFunc("/{id}", GetLocationImage).Methods("GET")

	locationsRoute := router.PathPrefix("/locations").Subrouter()
	locationsRoute.Handle("", user.VerifyJwtToken(http.HandlerFunc(controllerGetLocations))).Methods("GET")

}

func ResolveLocationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uri, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), locationContextKey, uri)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

func controllerGetLocation(w http.ResponseWriter, r *http.Request) {
	locationID := uint(r.Context().Value(locationContextKey).(int64))
	userID, ok := r.Context().Value(user.UserIdContextKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("Retreived user: %d and location: %d\n", userID, locationID)
}

func controllerDeleteLocation(w http.ResponseWriter, r *http.Request) {
	locationID := uint(r.Context().Value(locationContextKey).(int64))
	userID, ok := r.Context().Value(user.UserIdContextKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("Retreived user: %d and location: %d\n", userID, locationID)
}

func controllerGetLocations(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(user.UserIdContextKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var locations []Location
	config.DBConnection.Where("user_id = ?", userID).Find(&locations)

	w.Header().Set("Content-Type", "application/json")
	payload, _ := json.Marshal(locations)
	_, _ = w.Write(payload)
}
