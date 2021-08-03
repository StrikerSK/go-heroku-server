package location

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-heroku-server/api/src"
	"go-heroku-server/api/user"
	"net/http"
	"strconv"
)

const locationContextKey = "locationID"

func EnrichRouteWithLocation(router *mux.Router) {

	subroute := router.PathPrefix("/location").Subrouter()
	subroute.Handle("/add", user.VerifyJwtToken(http.HandlerFunc(controllerAddLocation))).Methods("POST")
	subroute.Handle("/{id}", user.VerifyJwtToken(ResolveLocationID(http.HandlerFunc(controllerUpdateLocation)))).Methods("PUT")
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
		uri, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			src.NewErrorResponse(http.StatusBadRequest, err).WriteResponse(w)
			return
		}
		ctx := context.WithValue(r.Context(), locationContextKey, uri)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func controllerAddLocation(w http.ResponseWriter, r *http.Request) {
	var res src.IResponse

	userID, res := user.ResolveUserContext(r.Context())
	if res != nil {
		res.WriteResponse(w)
		return
	}

	var location Location
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		res = src.NewErrorResponse(http.StatusInternalServerError, err)
		res.WriteResponse(w)
		return
	}

	addLocation(userID, location)
}

func controllerUpdateLocation(w http.ResponseWriter, r *http.Request) {
	var res src.IResponse

	locationID := resolveLocationContext(r.Context())
	userID, res := user.ResolveUserContext(r.Context())
	if res != nil {
		res.WriteResponse(w)
		return
	}

	var location Location
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		res = src.NewErrorResponse(http.StatusInternalServerError, err)
		res.WriteResponse(w)
		return
	}

	res = editLocation(userID, locationID, location)
	res.WriteResponse(w)
}

func controllerGetLocation(w http.ResponseWriter, r *http.Request) {
	var res src.IResponse

	userID, res := user.ResolveUserContext(r.Context())
	if res != nil {
		res.WriteResponse(w)
		return
	}

	res = retrieveLocation(userID, resolveLocationContext(r.Context()))
	res.WriteResponse(w)
}

func controllerDeleteLocation(w http.ResponseWriter, r *http.Request) {
	locationID := resolveLocationContext(r.Context())
	userID, err := user.ResolveUserContext(r.Context())
	if err != nil {
		err.WriteResponse(w)
		return
	}

	if res := deleteLocation(userID, locationID); err != nil {
		res.WriteResponse(w)
		return
	}
}

func controllerGetLocations(w http.ResponseWriter, r *http.Request) {
	var res src.IResponse
	userID, res := user.ResolveUserContext(r.Context())
	if res != nil {
		res.WriteResponse(w)
		return
	}

	res = getAllLocations(userID)
	res.WriteResponse(w)
}

func resolveLocationContext(context context.Context) uint {
	return uint(context.Value(locationContextKey).(uint64))
}
