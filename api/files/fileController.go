package files

import (
	"context"
	"github.com/gorilla/mux"
	"go-heroku-server/api/user"
	"go-heroku-server/config"
	"net/http"
	"strconv"
)

func EnrichRouteWithFile(router *mux.Router) {

	config.DBConnection.AutoMigrate(&File{})

	subroute := router.PathPrefix("/file").Subrouter()
	subroute.Handle("/upload", user.VerifyJwtToken(http.HandlerFunc(uploadFile))).Methods("POST")
	subroute.Handle("/{id}", ResolveFileID(user.VerifyJwtToken(http.HandlerFunc(readFile)))).Methods("GET")
	subroute.Handle("/{id}", ResolveFileID(user.VerifyJwtToken(http.HandlerFunc(removeFile)))).Methods("DELETE")

	filesSubroute := router.PathPrefix("/files").Subrouter()
	filesSubroute.Handle("/", user.VerifyJwtToken(http.HandlerFunc(getFileList))).Methods("GET")

}

func ResolveFileID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uri, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "file_ID", uri)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
