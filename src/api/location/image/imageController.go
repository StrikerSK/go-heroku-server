package image

import (
	"github.com/gorilla/mux"
	"go-heroku-server/config"
	"go-heroku-server/src/responses"
	"log"
	"net/http"
	"strconv"
)

func EnrichRouteWithImages(router *mux.Router) {
	config.InitializeType("LocationImage", &LocationImage{})

	imageRoute := router.PathPrefix("/image").Subrouter()
	imageRoute.HandleFunc("/{imageId}", cReadLocationImage).Methods(http.MethodGet)
}

func cReadLocationImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID, err := strconv.ParseInt(vars["imageId"], 10, 64)
	if err != nil {
		log.Printf("Location image [%d] read: %s\n", imageID, err)
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	sReadLocationImage(imageID).WriteResponse(w)
	return
}
