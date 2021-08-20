package image

import (
	"github.com/gorilla/mux"
	"net/http"
)

func EnrichRouteWithImages(router *mux.Router) {
	imageRoute := router.PathPrefix("/image").Subrouter()
	imageRoute.HandleFunc("/{imageId}", getLocationImage).Methods(http.MethodGet)
}
