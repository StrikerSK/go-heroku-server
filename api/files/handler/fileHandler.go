package fileHandlers

import (
	"fmt"
	"github.com/gorilla/mux"
	fileDomains "go-heroku-server/api/files/domain"
	filePorts "go-heroku-server/api/files/ports"
	"go-heroku-server/api/src/errors"
	"go-heroku-server/api/src/responses"
	userHandlers "go-heroku-server/api/user/handler"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FileHandler struct {
	fileService     filePorts.IFileService
	userMiddleware  userHandlers.UserAuthMiddleware
	responseService responses.ResponseService
}

func NewFileHandler(service filePorts.IFileService, userMiddleware userHandlers.UserAuthMiddleware, responseService responses.ResponseService) FileHandler {
	return FileHandler{
		fileService:     service,
		userMiddleware:  userMiddleware,
		responseService: responseService,
	}
}

func (h FileHandler) EnrichRouter(router *mux.Router) {
	fileRoute := router.PathPrefix("/file").Subrouter()
	fileRoute.Handle("/upload", h.userMiddleware.VerifyToken(http.HandlerFunc(h.createFile))).Methods(http.MethodPost)
	fileRoute.Handle("/{id}", h.userMiddleware.VerifyToken(http.HandlerFunc(h.readFile))).Methods(http.MethodGet)
	fileRoute.Handle("/{id}", h.userMiddleware.VerifyToken(http.HandlerFunc(h.deleteFile))).Methods(http.MethodDelete)

	filesRoute := router.PathPrefix("/files").Subrouter()
	filesRoute.Handle("/", h.userMiddleware.VerifyToken(http.HandlerFunc(h.readFiles))).Methods(http.MethodGet)

}

func (h FileHandler) createFile(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUsernameFromContext(r.Context())
	if err != nil {
		log.Printf("Controller file upload: %s\n", err.Error())
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("Controller file upload: %s\n", err.Error())
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	//Processing of received file metadata
	if err = r.ParseForm(); err != nil {
		log.Printf("Controller file upload: %s\n", err.Error())
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("File add: %s\n", err.Error())
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	contentType := fileHeader.Header.Get("Content-Type")

	resolvedFile := fileDomains.FileEntity{
		Username:   username,
		FileName:   fileHeader.Filename,
		FileType:   contentType,
		FileData:   fileBytes,
		FileSize:   h.getFileSize(fileHeader.Size),
		CreateDate: time.Now(),
	}

	id, err := h.fileService.CreateFile(resolvedFile)
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	} else {
		log.Printf("File create: success\n")
		h.responseService.CreateResponse(map[string]interface{}{"id": id}).WriteResponse(w)
		return
	}
}

func (h FileHandler) readFile(w http.ResponseWriter, r *http.Request) {
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

	res := responses.CreateResponse(http.StatusOK, persistedFile.FileData)
	res.SetHeaders(responseMap)
	res.WriteResponse(w)

	return
}

func (h FileHandler) deleteFile(w http.ResponseWriter, r *http.Request) {
	username, err := h.userMiddleware.GetUsernameFromContext(r.Context())
	if err != nil {
		h.responseService.CreateResponse(err).WriteResponse(w)
		return
	}

	fileID, err := h.resolveFileIdentificationContext(r)
	if err != nil {
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

func (h FileHandler) readFiles(w http.ResponseWriter, r *http.Request) {
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

func (FileHandler) resolveFileIdentificationContext(r *http.Request) (uint, error) {
	tmpVar := mux.Vars(r)["id"]
	uri, err := strconv.ParseUint(tmpVar, 10, 64)
	if err != nil {
		return 0, errors.NewParseError(err.Error())
	}
	return uint(uri), err
}

//Resolve ideal file size up to MegaBytes
func (FileHandler) getFileSize(fileSize int64) (outputSize string) {
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
