package files

import (
	"context"
	"github.com/gorilla/mux"
	"go-heroku-server/api/src/responses"
	"go-heroku-server/api/user"
	"go-heroku-server/config"
	"log"
	"net/http"
	"strconv"
)

const fileContextName = "file_ID"

func EnrichRouteWithFile(router *mux.Router) {

	config.DBConnection.AutoMigrate(&File{})

	fileRoute := router.PathPrefix("/file").Subrouter()
	fileRoute.Handle("/upload", user.VerifyJwtToken(http.HandlerFunc(controllerUploadFile))).Methods(http.MethodPost)
	fileRoute.Handle("/{id}", resolveFileID(user.VerifyJwtToken(http.HandlerFunc(controllerReadFile)))).Methods(http.MethodGet)
	fileRoute.Handle("/{id}", resolveFileID(user.VerifyJwtToken(http.HandlerFunc(controllerRemoveFile)))).Methods(http.MethodDelete)

	filesRoute := router.PathPrefix("/files").Subrouter()
	filesRoute.Handle("/", user.VerifyJwtToken(http.HandlerFunc(controllerGetFileList))).Methods(http.MethodGet)

}

func resolveFileID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uri, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.Printf("Request file resolving: %s\n", err.Error())
			responses.NewEmptyResponse(http.StatusBadRequest).WriteResponse(w)
			return
		}
		ctx := context.WithValue(r.Context(), fileContextName, uri)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func controllerUploadFile(w http.ResponseWriter, r *http.Request) {
	userID, res := user.ResolveUserContext(r.Context())
	if res != nil {
		res.WriteResponse(w)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("Controller file upload: %s\n", err.Error())
		responses.NewEmptyResponse(http.StatusBadRequest).WriteResponse(w)
		return
	}

	//Processing of received file metadata
	if err = r.ParseForm(); err != nil {
		log.Printf("Controller file upload: %s\n", err.Error())
		responses.NewEmptyResponse(http.StatusInternalServerError).WriteResponse(w)
		return
	}

	//r.FormValue("file") -> Receives metadata from sent file
	//Process of JSON to items (Attention! Struct form must be predefined)
	//var token user.Token
	//if err = json.Unmarshal([]byte(r.FormValue("file")), &token); err != nil {
	//	log.Print(err)
	//}

	uploadFile(file, fileHeader, userID).WriteResponse(w)
}

func controllerReadFile(w http.ResponseWriter, r *http.Request) {
	userID, _ := user.ResolveUserContext(r.Context())
	fileID := resolveFileContext(r.Context())
	readFile(userID, fileID).WriteResponse(w)
	return
}

func controllerRemoveFile(w http.ResponseWriter, r *http.Request) {
	userID, _ := user.ResolveUserContext(r.Context())
	removeFile(userID, resolveFileContext(r.Context())).WriteResponse(w)
}

func controllerGetFileList(w http.ResponseWriter, r *http.Request) {
	if userID, res := user.ResolveUserContext(r.Context()); res != nil {
		res.WriteResponse(w)
	} else {
		getFileList(userID).WriteResponse(w)
	}
}

func resolveFileContext(context context.Context) uint {
	return uint(context.Value(fileContextName).(int64))
}
