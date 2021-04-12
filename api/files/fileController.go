package files

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-heroku-server/api/user"
	"go-heroku-server/config"
	"log"
	"net/http"
	"strconv"
)

func EnrichRouteWithFile(router *mux.Router) {

	config.DBConnection.AutoMigrate(&File{})

	subroute := router.PathPrefix("/file").Subrouter()
	subroute.Handle("/upload", user.VerifyJwtToken(http.HandlerFunc(controllerUploadFile))).Methods("POST")
	subroute.Handle("/{id}", ResolveFileID(user.VerifyJwtToken(http.HandlerFunc(controllerReadFile)))).Methods("GET")
	subroute.Handle("/{id}", ResolveFileID(user.VerifyJwtToken(http.HandlerFunc(controllerRemoveFile)))).Methods("DELETE")

	filesSubroute := router.PathPrefix("/files").Subrouter()
	filesSubroute.Handle("/", user.VerifyJwtToken(http.HandlerFunc(controllerGetFileList))).Methods("GET")

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

func controllerUploadFile(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return
	}

	//Processing of received file metadata
	if err = r.ParseForm(); err != nil {
		log.Print(err)
	}

	//r.FormValue("file") -> Receives metadata from sent file
	//Process of JSON to items (Attention! Struct form must be predefined)
	//var token user.Token
	//if err = json.Unmarshal([]byte(r.FormValue("file")), &token); err != nil {
	//	log.Print(err)
	//}

	uploadFile(file, fileHeader, r.Context().Value(user.UserIdContextKey).(uint))
}

func controllerReadFile(w http.ResponseWriter, r *http.Request) {
	fileID := uint(r.Context().Value("file_ID").(int64))
	userID := r.Context().Value(user.UserIdContextKey).(uint)

	if file, requestError := readFile(userID, fileID); requestError != nil {
		w.WriteHeader(requestError.StatusCode)
		return
	} else {
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, Content-Length, X-Content-Transfer-Id")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Disposition", "attachment; filename="+file.FileName)
		w.Header().Set("Content-Type", file.FileType)
		_, _ = w.Write(file.FileData)
	}
}

func controllerRemoveFile(w http.ResponseWriter, r *http.Request) {
	fileID := uint(r.Context().Value("file_ID").(int64))
	userID := r.Context().Value(user.UserIdContextKey).(uint)
	if requestError := removeFile(userID, fileID); requestError != nil {
		w.WriteHeader(requestError.StatusCode)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func controllerGetFileList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_ = json.NewEncoder(w).Encode(getFileList(r.Context().Value(user.UserIdContextKey).(uint)))
}
