package fileHandlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	fileDomains "go-heroku-server/api/files/domain"
	filePorts "go-heroku-server/api/files/ports"
	"go-heroku-server/api/src/responses"
	userHandlers "go-heroku-server/api/user/handler"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const fileContextName = "file_ID"

type FileHandler struct {
	fileService    filePorts.IFileService
	userMiddleware userHandlers.UserAuthMiddleware
}

func NewFileHandler(service filePorts.IFileService, userMiddleware userHandlers.UserAuthMiddleware) FileHandler {
	return FileHandler{
		fileService:    service,
		userMiddleware: userMiddleware,
	}
}

func (h FileHandler) EnrichRouter(router *mux.Router) {
	fileRoute := router.PathPrefix("/file").Subrouter()
	fileRoute.Handle("/upload", h.userMiddleware.VerifyToken(http.HandlerFunc(h.createFile))).Methods(http.MethodPost)
	fileRoute.Handle("/{id}", resolveFileID(h.userMiddleware.VerifyToken(http.HandlerFunc(h.readFile)))).Methods(http.MethodGet)
	fileRoute.Handle("/{id}", resolveFileID(h.userMiddleware.VerifyToken(http.HandlerFunc(h.deleteFile)))).Methods(http.MethodDelete)

	filesRoute := router.PathPrefix("/files").Subrouter()
	filesRoute.Handle("/", h.userMiddleware.VerifyToken(http.HandlerFunc(h.readFiles))).Methods(http.MethodGet)

}

func resolveFileID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uri, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.Printf("Request file resolving: %s\n", err.Error())
			responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
			return
		}
		ctx := context.WithValue(r.Context(), fileContextName, uri)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h FileHandler) createFile(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUserFromContext(r.Context())
	if err != nil {
		log.Printf("Controller file upload: %s\n", err.Error())
		responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("Controller file upload: %s\n", err.Error())
		responses.CreateResponse(http.StatusBadRequest, nil).WriteResponse(w)
		return
	}

	//Processing of received file metadata
	if err = r.ParseForm(); err != nil {
		log.Printf("Controller file upload: %s\n", err.Error())
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("File add: %s\n", err.Error())
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")

	resolvedFile := fileDomains.FileEntity{
		Username:   username,
		FileName:   fileHeader.Filename,
		FileType:   contentType,
		FileData:   fileBytes,
		FileSize:   getFileSize(fileHeader.Size),
		CreateDate: time.Now(),
	}

	if err = h.fileService.UploadFile(resolvedFile); err != nil {
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	} else {
		log.Printf("File create: success\n")
		return
	}
}

func (h FileHandler) readFile(w http.ResponseWriter, r *http.Request) {
	username, _ := h.userMiddleware.GetUserFromContext(r.Context())
	fileID := resolveFileContext(r.Context())

	var persistedFile, err = h.fileService.ReadFile(fileID, username)
	if err != nil {
		log.Printf("File [%d] read: %v\n", fileID, err)
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	responseMap := map[string]string{
		"Access-Control-Expose-Headers": "Content-Disposition, Content-Length, X-Content-Transfer-Id",
		"Access-Control-Allow-Origin":   "*",
		"Content-Disposition":           "attachment; filename=" + persistedFile.FileName,
		"Content-Type":                  persistedFile.FileType,
	}

	res := responses.CreateResponse(http.StatusOK, persistedFile.FileData)
	res.SetHeaders(responseMap)
	res.WriteResponse(w)

	return
}

func (h FileHandler) deleteFile(w http.ResponseWriter, r *http.Request) {
	username, _ := h.userMiddleware.GetUserFromContext(r.Context())
	fileID := resolveFileContext(r.Context())

	if err := h.fileService.DeleteFile(fileID, username); err != nil {
		log.Printf("File [%d] delete: %s\n", fileID, err.Error())
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	responses.CreateResponse(http.StatusOK, nil).WriteResponse(w)
	return
}

func (h FileHandler) readFiles(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUserFromContext(r.Context())
	if err != nil {
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	}

	files, err := h.fileService.ReadFiles(username)
	if err != nil {
		log.Printf("Files read error occured: %v\n", err.Error())
		responses.CreateResponse(http.StatusInternalServerError, nil).WriteResponse(w)
		return
	} else {
		responses.CreateResponse(http.StatusOK, files).WriteResponse(w)
		return
	}
}

func resolveFileContext(context context.Context) uint {
	return uint(context.Value(fileContextName).(int64))
}

//Resolve ideal file size up to MegaBytes
func getFileSize(fileSize int64) (outputSize string) {
	switch {
	case fileSize < 1024:
		outputSize = fmt.Sprintf("%d B", fileSize)
		break
	case fileSize < 1048576:
		fileSize = fileSize / 1024
		outputSize = fmt.Sprintf("%d kB", fileSize)
		break
	default:
		fileSize = fileSize / 1048576
		outputSize = fmt.Sprintf("%d MB", fileSize)
		break
	}
	return
}
