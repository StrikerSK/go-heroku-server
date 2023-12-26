package fileHandlers

import (
	"github.com/gorilla/mux"
	fileUtils "go-heroku-server/api/files/utils"
	fileDomains "go-heroku-server/api/files/v2/domain"
	filePorts "go-heroku-server/api/files/v2/ports"
	"go-heroku-server/api/src/errors"
	"go-heroku-server/api/src/responses"
	userHandlers "go-heroku-server/api/user/handler"
	"io"
	"log"
	"net/http"
	"time"
)

type MuxFileHandlerV2 struct {
	fileService     filePorts.IFileServiceV2
	userMiddleware  userHandlers.UserAuthMiddleware
	responseService responses.ResponseFactory
}

func NewMuxFileHandler(service filePorts.IFileServiceV2, userMiddleware userHandlers.UserAuthMiddleware, responseService responses.ResponseFactory) MuxFileHandlerV2 {
	return MuxFileHandlerV2{
		fileService:     service,
		userMiddleware:  userMiddleware,
		responseService: responseService,
	}
}

func (h MuxFileHandlerV2) EnrichRouter(router *mux.Router) {
	fileRoute := router.PathPrefix("/file").Subrouter()
	fileRoute.Use(h.userMiddleware.VerifyToken)
	fileRoute.Handle("/upload", http.HandlerFunc(h.createFile)).Methods(http.MethodPost)
	fileRoute.Handle("/{id}", http.HandlerFunc(h.downloadFile)).Methods(http.MethodGet)
	fileRoute.Handle("/{id}", http.HandlerFunc(h.deleteFile)).Methods(http.MethodDelete)

	headerRoute := router.PathPrefix("/header").Subrouter()
	headerRoute.Use(h.userMiddleware.VerifyToken)
	headerRoute.Handle("/{id}", http.HandlerFunc(h.readMetadata)).Methods(http.MethodGet)
	headerRoute.Handle("", http.HandlerFunc(h.readFiles)).Methods(http.MethodGet)
}

func (h MuxFileHandlerV2) createFile(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUsername(r.Context())

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

	metadata := fileDomains.FileMetadataV2{
		Username:   username,
		FileName:   attachmentName,
		FileType:   contentType,
		FileSize:   fileUtils.GetFileSize(fileSize),
		CreateDate: time.Now(),
	}

	fileData := fileDomains.FileEntityV2{
		FileData: fileBytes,
	}

	fileObject := fileDomains.FileObjectV2{
		FileEntityV2:   fileData,
		FileMetadataV2: metadata,
	}

	id, err := h.fileService.CreateFile(fileObject)
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	} else {
		response := map[string]interface{}{"id": id}
		h.responseService.CreateResponse(response).WriteResponse(w)
		return
	}
}

func (h MuxFileHandlerV2) downloadFile(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUsername(r.Context())
	if err != nil {
		log.Printf("File read: %v\n", err)
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	fileID := h.resolveFileIdentificationContext(r)
	persistedFile, err := h.fileService.DownloadFile(fileID, username)
	if err != nil {
		log.Printf("File [%s] read: %v\n", fileID, err)
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	responseMap := map[string]string{
		"Access-Control-Expose-Headers": "Content-Disposition, Content-Length, X-Content-Transfer-Id",
		"Access-Control-Allow-Origin":   "*",
		"Content-Disposition":           "attachment; filename=" + persistedFile.FileName,
		"Content-Type":                  persistedFile.FileType,
	}

	res := h.responseService.CreateResponse(persistedFile.FileEntityV2.FileData)
	res.SetHeaders(responseMap)
	res.WriteResponse(w)
	return
}

func (h MuxFileHandlerV2) readMetadata(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUsername(r.Context())
	if err != nil {
		log.Printf("File read: %v\n", err)
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	fileID := h.resolveFileIdentificationContext(r)
	persistedFile, err := h.fileService.ReadMetadata(fileID, username)
	if err != nil {
		log.Printf("File [%s] read: %v\n", fileID, err)
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	res := h.responseService.CreateResponse(persistedFile)
	res.WriteResponse(w)
	return
}

func (h MuxFileHandlerV2) deleteFile(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUsername(r.Context())
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	fileID := h.resolveFileIdentificationContext(r)

	if err = h.fileService.RemoveFile(fileID, username); err != nil {
		log.Printf("File [%s] delete: %s\n", fileID, err.Error())
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	h.responseService.CreateResponse(nil).WriteResponse(w)
	return
}

func (h MuxFileHandlerV2) readFiles(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUsername(r.Context())
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

func (MuxFileHandlerV2) resolveFileIdentificationContext(r *http.Request) string {
	return mux.Vars(r)["id"]
}
