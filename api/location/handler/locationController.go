package locationHandlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-heroku-server/api/location/domain"
	"go-heroku-server/api/location/image"
	locationPorts "go-heroku-server/api/location/ports"
	"go-heroku-server/api/src/responses"
	userHandlers "go-heroku-server/api/user/handler"
	"log"
	"net/http"
	"strconv"
)

const locationContextKey = "locationID"

type LocationHandler struct {
	locationService locationPorts.ILocationService
	userMiddleware  userHandlers.UserAuthMiddleware
}

func NewLocationHandler(locationService locationPorts.ILocationService, userMiddleware userHandlers.UserAuthMiddleware) LocationHandler {
	return LocationHandler{
		locationService: locationService,
		userMiddleware:  userMiddleware,
	}
}

func (h LocationHandler) EnrichRouter(router *mux.Router) {
	locationRoute := router.PathPrefix("/location").Subrouter()
	locationRoute.Handle("", h.userMiddleware.VerifyToken(http.HandlerFunc(h.createLocation))).Methods(http.MethodPost)

	locationSubRoute := locationRoute.PathPrefix("/{id}").Subrouter()
	locationSubRoute.Handle("", h.userMiddleware.VerifyToken(ResolveLocationID(http.HandlerFunc(h.updateLocation)))).Methods(http.MethodPut)
	locationSubRoute.Handle("", h.userMiddleware.VerifyToken(ResolveLocationID(http.HandlerFunc(h.readLocation)))).Methods(http.MethodGet)
	locationSubRoute.Handle("", h.userMiddleware.VerifyToken(ResolveLocationID(http.HandlerFunc(h.deleteLocation)))).Methods(http.MethodDelete)

	locationsRoute := router.PathPrefix("/locations").Subrouter()
	locationsRoute.Handle("", h.userMiddleware.VerifyToken(http.HandlerFunc(h.readLocations))).Methods(http.MethodGet)

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

func (h LocationHandler) createLocation(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUserFromContext(r.Context())
	if err != nil {
		log.Printf("Controller location add: %v\n", err)
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	var location locationDomains.UserLocationEntity
	if err = json.NewDecoder(r.Body).Decode(&location); err != nil {
		log.Printf("Controller location add: %v\n", err)
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	location.Username = username

	err = h.locationService.CreateLocation(location)
	if err != nil {
		log.Printf("Controller location add: %v\n", err)
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	} else {
		responses.CreateResponse(http.StatusOK, nil).WriteResponse(w)
		return
	}
}

func (h LocationHandler) readLocation(w http.ResponseWriter, r *http.Request) {
	locationID := resolveLocationContext(r.Context())
	username, err := h.userMiddleware.GetUserFromContext(r.Context())
	if err != nil {
		log.Printf("Read location error: %v\n", err)
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	location, err := h.locationService.GetLocation(locationID, username)
	if err != nil {
		log.Printf("Location [%d] read: %s", locationID, err.Error())
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	responses.CreateResponse(http.StatusOK, location).WriteResponse(w)
	return
}

func (h LocationHandler) readLocations(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUserFromContext(r.Context())

	if err != nil {
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	locations, err := h.locationService.ReadLocations(username)
	if err != nil {
		log.Printf("Locations reading: %s", err.Error())
		responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
		return
	} else {
		responses.CreateResponse(http.StatusOK, locations).WriteResponse(w)
		return
	}
}

func (h LocationHandler) updateLocation(w http.ResponseWriter, r *http.Request) {
	locationID := resolveLocationContext(r.Context())
	username, err := h.userMiddleware.GetUserFromContext(r.Context())
	if err != nil {
		log.Printf("Controller location add: %v\n", err)
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	var location locationDomains.UserLocationEntity
	if err = json.NewDecoder(r.Body).Decode(&location); err != nil {
		log.Printf("Location [%d] edit: %s\n", locationID, err.Error())
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	location.Id = locationID
	location.Username = username

	err = h.locationService.UpdateLocation(location)
	if err != nil {
		log.Printf("Location [%d] edit: %s\n", locationID, err.Error())
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	} else {
		responses.CreateResponse(http.StatusOK, nil).WriteResponse(w)
		return
	}
}

func (h LocationHandler) deleteLocation(w http.ResponseWriter, r *http.Request) {
	locationID := resolveLocationContext(r.Context())

	username, err := h.userMiddleware.GetUserFromContext(r.Context())
	if err != nil {
		log.Printf("Location delete: %v", err)
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	err = h.locationService.DeleteLocation(locationID, username)
	if err != nil {
		log.Printf("Location [%d] delete: %v", locationID, err)
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	responses.CreateResponse(http.StatusOK, nil).WriteResponse(w)
	return
}

func resolveLocationContext(context context.Context) uint {
	return uint(context.Value(locationContextKey).(uint64))
}
