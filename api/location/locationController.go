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

	locationRoute := router.PathPrefix("/location").Subrouter()
	locationRoute.Handle("/add", user.VerifyJwtToken(http.HandlerFunc(controllerAddLocation))).Methods(http.MethodPost)
	locationRoute.Handle("/{id}", user.VerifyJwtToken(ResolveLocationID(http.HandlerFunc(controllerUpdateLocation)))).Methods(http.MethodPut)
	locationRoute.Handle("/{id}", user.VerifyJwtToken(ResolveLocationID(http.HandlerFunc(controllerGetLocation)))).Methods(http.MethodGet)
	locationRoute.Handle("/{id}", user.VerifyJwtToken(ResolveLocationID(http.HandlerFunc(controllerDeleteLocation)))).Methods(http.MethodDelete)

	imageSubRoute := locationRoute.PathPrefix("/image").Subrouter()
	imageSubRoute.HandleFunc("/{id}", GetLocationImage).Methods(http.MethodGet)

	locationsRoute := router.PathPrefix("/locations").Subrouter()
	locationsRoute.Handle("", user.VerifyJwtToken(http.HandlerFunc(controllerGetLocations))).Methods(http.MethodGet)

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
	var res src.IResponse
	locationID := resolveLocationContext(r.Context())
	userID, res := user.ResolveUserContext(r.Context())
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
