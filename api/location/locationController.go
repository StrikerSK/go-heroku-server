package location

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-heroku-server/api/location/image"
	"go-heroku-server/api/src/responses"
	"go-heroku-server/config"
	"log"
	"net/http"
	"strconv"
)

const locationContextKey = "locationID"

func EnrichRouteWithLocation(router *mux.Router) {

	config.InitializeType("Location", &Location{})

	locationRoute := router.PathPrefix("/location").Subrouter()
	locationRoute.Handle("", handler.VerifyJwtToken(http.HandlerFunc(controllerAddLocation))).Methods(http.MethodPost)

	locationSubRoute := locationRoute.PathPrefix("/{id}").Subrouter()
	locationSubRoute.Handle("", handler.VerifyJwtToken(ResolveLocationID(http.HandlerFunc(controllerUpdateLocation)))).Methods(http.MethodPut)
	locationSubRoute.Handle("", handler.VerifyJwtToken(ResolveLocationID(http.HandlerFunc(controllerGetLocation)))).Methods(http.MethodGet)
	locationSubRoute.Handle("", handler.VerifyJwtToken(ResolveLocationID(http.HandlerFunc(controllerDeleteLocation)))).Methods(http.MethodDelete)

	locationsRoute := router.PathPrefix("/locations").Subrouter()
	locationsRoute.Handle("", handler.VerifyJwtToken(http.HandlerFunc(controllerGetLocations))).Methods(http.MethodGet)

	image.EnrichRouteWithImages(locationSubRoute)
}

func ResolveLocationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uri, err := strconv.ParseUint(vars["id"], 10, 64)
		if err != nil {
			log.Printf("Resolving location: %s\n", err.Error())
			responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
			return
		}
		ctx := context.WithValue(r.Context(), locationContextKey, uri)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func controllerAddLocation(w http.ResponseWriter, r *http.Request) {
	var res responses.IResponse

	userID, res := handler.ResolveUserContext(r.Context())
	if res != nil {
		res.WriteResponse(w)
		return
	}

	var location UserLocation
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		log.Printf("Controller location add: %s\n", err.Error())
		res = responses.CreateResponse(http.StatusInternalServerError, nil)
		res.WriteResponse(w)
		return
	}

	addLocation(userID, location)
}

func controllerUpdateLocation(w http.ResponseWriter, r *http.Request) {
	var res responses.IResponse

	locationID := resolveLocationContext(r.Context())
	userID, res := handler.ResolveUserContext(r.Context())
	if res != nil {
		res.WriteResponse(w)
		return
	}

	var location UserLocation
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		log.Printf("Location [%d] edit (Controller): %s\n", locationID, err.Error())
		res = responses.CreateResponse(http.StatusInternalServerError, nil)
		res.WriteResponse(w)
		return
	}

	res = editLocation(userID, locationID, location)
	res.WriteResponse(w)
}

func controllerGetLocation(w http.ResponseWriter, r *http.Request) {
	var res responses.IResponse

	userID, res := handler.ResolveUserContext(r.Context())
	if res != nil {
		res.WriteResponse(w)
		return
	}

	res = retrieveLocation(userID, resolveLocationContext(r.Context()))
	res.WriteResponse(w)
}

func controllerDeleteLocation(w http.ResponseWriter, r *http.Request) {
	var res responses.IResponse
	locationID := resolveLocationContext(r.Context())
	userID, res := handler.ResolveUserContext(r.Context())
	if res != nil {
		res.WriteResponse(w)
		return
	}

	if res = deleteLocation(userID, locationID); res != nil {
		res.WriteResponse(w)
		return
	}
}

func controllerGetLocations(w http.ResponseWriter, r *http.Request) {
	var res responses.IResponse
	userID, res := handler.ResolveUserContext(r.Context())
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
