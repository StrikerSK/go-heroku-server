package fileHandlers

import (
	"github.com/gorilla/mux"
	fileDomains "go-heroku-server/api/files/domain"
	filePorts "go-heroku-server/api/files/ports"
	fileUtils "go-heroku-server/api/files/utils"
	"go-heroku-server/api/src/errors"
	"go-heroku-server/api/src/responses"
	userHandlers "go-heroku-server/api/user/handler"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type MuxFileHandler struct {
	fileService     filePorts.IFileService
	userMiddleware  userHandlers.UserAuthMiddleware
	responseService responses.ResponseFactory
}

func NewMuxFileHandler(service filePorts.IFileService, userMiddleware userHandlers.UserAuthMiddleware, responseService responses.ResponseFactory) MuxFileHandler {
	return MuxFileHandler{
		fileService:     service,
		userMiddleware:  userMiddleware,
		responseService: responseService,
	}
}

func (h MuxFileHandler) EnrichRouter(router *mux.Router) {
	fileRoute := router.PathPrefix("/file").Subrouter()
	fileRoute.Use(h.userMiddleware.VerifyToken)
	fileRoute.Handle("/upload", http.HandlerFunc(h.createFile)).Methods(http.MethodPost)
	fileRoute.Handle("/{id}", http.HandlerFunc(h.readFile)).Methods(http.MethodGet)
	fileRoute.Handle("/{id}", http.HandlerFunc(h.deleteFile)).Methods(http.MethodDelete)

	filesRoute := router.PathPrefix("/files").Subrouter()
	filesRoute.Use(h.userMiddleware.VerifyToken)
	filesRoute.Handle("/", http.HandlerFunc(h.readFiles)).Methods(http.MethodGet)

}

func (h MuxFileHandler) createFile(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUsernameFromContext(r.Context())

	if err != nil {
		log.Printf("Controller file upload: %s\n", err.Error())
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	attachmentName := r.URL.Query().Get("name")
	if attachmentName == "" {
		err := errors.NewBadRequestError("Attachment's name not provided")
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	fileBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("File add: %s\n", err.Error())
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	fileSize := int64(len(fileBytes))
	if fileSize <= 0 {
		err := errors.NewBadRequestError("Attachment should not be empty")
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	contentType := http.DetectContentType(fileBytes)

	resolvedFile := fileDomains.FileEntity{
		Username:   username,
		FileName:   attachmentName,
		FileType:   contentType,
		FileData:   fileBytes,
		FileSize:   fileUtils.GetFileSize(fileSize),
		CreateDate: time.Now(),
	}

	id, err := h.fileService.CreateFile(resolvedFile)
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	} else {
		response := map[string]interface{}{"id": id}
		h.responseService.CreateResponse(response).WriteResponse(w)
		return
	}
}

func (h MuxFileHandler) readFile(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUsernameFromContext(r.Context())
	if err != nil {
		log.Printf("File read: %v\n", err)
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	fileID, err := h.resolveFileIdentificationContext(r)
	if err != nil {
		log.Printf("File [%d] read: %v\n", fileID, err)
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	persistedFile, err := h.fileService.ReadFile(fileID, username)
	if err != nil {
		log.Printf("File [%d] read: %v\n", fileID, err)
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	responseMap := map[string]string{
		"Access-Control-Expose-Headers": "Content-Disposition, Content-Length, X-Content-Transfer-Id",
		"Access-Control-Allow-Origin":   "*",
		"Content-Disposition":           "attachment; filename=" + persistedFile.FileName,
		"Content-Type":                  persistedFile.FileType,
	}

	res := h.responseService.CreateResponse(persistedFile.FileData)
	res.SetHeaders(responseMap)
	res.WriteResponse(w)
	return
}

func (h MuxFileHandler) deleteFile(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUsernameFromContext(r.Context())
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	fileID, err := h.resolveFileIdentificationContext(r)
	if err != nil {
		log.Printf("File [%d] delete: %s\n", fileID, err.Error())
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	if err = h.fileService.DeleteFile(fileID, username); err != nil {
		log.Printf("File [%d] delete: %s\n", fileID, err.Error())
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	h.responseService.CreateResponse(nil).WriteResponse(w)
	return
}

func (h MuxFileHandler) readFiles(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUsernameFromContext(r.Context())
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	files, err := h.fileService.ReadFiles(username)
	if err != nil {
		log.Printf("Files read error occured: %v\n", err.Error())
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	} else {
		h.responseService.CreateResponse(files).WriteResponse(w)
		return
	}
}

func (MuxFileHandler) resolveFileIdentificationContext(r *http.Request) (uint, error) {
	tmpVar := mux.Vars(r)["id"]
	uri, err := strconv.ParseUint(tmpVar, 10, 64)
	if err != nil {
		return 0, errors.NewParseError(err.Error())
	}
	return uint(uri), err
}
