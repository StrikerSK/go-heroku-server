package files

import (
	"context"
	"github.com/gorilla/mux"
	"go-heroku-server/api/src"
	"go-heroku-server/api/user"
	"go-heroku-server/config"
	"net/http"
	"strconv"
)

const fileContextName = "file_ID"

func EnrichRouteWithFile(router *mux.Router) {

	config.DBConnection.AutoMigrate(&File{})

	subroute := router.PathPrefix("/file").Subrouter()
	subroute.Handle("/upload", user.VerifyJwtToken(http.HandlerFunc(controllerUploadFile))).Methods("POST")
	subroute.Handle("/{id}", resolveFileID(user.VerifyJwtToken(http.HandlerFunc(controllerReadFile)))).Methods("GET")
	subroute.Handle("/{id}", resolveFileID(user.VerifyJwtToken(http.HandlerFunc(controllerRemoveFile)))).Methods("DELETE")

	filesSubroute := router.PathPrefix("/files").Subrouter()
	filesSubroute.Handle("/", user.VerifyJwtToken(http.HandlerFunc(controllerGetFileList))).Methods("GET")

}

func resolveFileID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uri, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			src.NewErrorResponse(http.StatusBadRequest, err).WriteResponse(w)
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
		src.NewErrorResponse(http.StatusBadRequest, err).WriteResponse(w)
		return
	}

	//Processing of received file metadata
	if err = r.ParseForm(); err != nil {
		src.NewErrorResponse(http.StatusInternalServerError, err).WriteResponse(w)
		return
	}

	//r.FormValue("file") -> Receives metadata from sent file
	//Process of JSON to items (Attention! Struct form must be predefined)
	//var token user.Token
	//if err = json.Unmarshal([]byte(r.FormValue("file")), &token); err != nil {
	//	log.Print(err)
	//}

	uploadFile(file, fileHeader, userID)
}

func controllerReadFile(w http.ResponseWriter, r *http.Request) {
	userID, _ := user.ResolveUserContext(r.Context())

	if file, err := readFile(userID, resolveFileContext(r.Context())); err != nil {
		err.WriteResponse(w)
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
