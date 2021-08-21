package files

import (
	"fmt"
	"go-heroku-server/api/src/responses"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

//Function stores files received from the Front-End
func uploadFile(file multipart.File, fileHeader *multipart.FileHeader, userID uint) responses.IResponse {
	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		log.Printf("File add: %s\n", err.Error())
		return responses.NewEmptyResponse(http.StatusInternalServerError)
	}

	contentType := fileHeader.Header.Get("Content-Type")

	resolvedFile := File{
		UserID:     userID,
		FileName:   fileHeader.Filename,
		FileType:   contentType,
		FileData:   fileBytes,
		FileSize:   getFileSize(fileHeader.Size),
		CreateDate: time.Now(),
	}

	createFile(resolvedFile)
	log.Printf("File create: success\n")
	return responses.NewEmptyResponse(http.StatusCreated)
}

//Function provides requested file to the client
func readFile(userID uint, fileID uint) responses.IResponse {
	var persistedFile, err = getFile(fileID)
	if err != nil {
		log.Printf("File [%d] read: %s\n", fileID, err.Error())
		return responses.NewEmptyResponse(http.StatusNotFound)
	}

	if persistedFile.UserID != userID {
		log.Printf("File [%d] read: access denied\n", fileID)
		return responses.NewResponse(http.StatusForbidden)
	}

	responseMap := map[string]string{
		"Access-Control-Expose-Headers": "Content-Disposition, Content-Length, X-Content-Transfer-Id",
		"Access-Control-Allow-Origin":   "*",
		"Content-Disposition":           "attachment; filename=" + persistedFile.FileName,
		"Content-Type":                  persistedFile.FileType,
	}

	log.Printf("File [%d] read: success\n", fileID)
	return responses.NewFileResponse(persistedFile.FileData, responseMap)
}

func getFileList(userID uint) responses.ResponseImpl {
	files := getAll(userID)
	for index := range files {
		fileName := files[index].FileName
		fileName = fileName[:strings.IndexByte(fileName, '.')]
		files[index].FileName = fileName
	}
	log.Printf("File listing: success\n")
	return responses.NewResponse(files)
}

//Deletion of file base on userID
func removeFile(userID, fileID uint) responses.IResponse {
	persistedFile, err := getFile(fileID)
	if err != nil {
		log.Printf("File [%d] delete: %s\n", fileID, err.Error())
		return responses.NewEmptyResponse(http.StatusOK)
	}

	if persistedFile.UserID != userID {
		log.Printf("File [%d] delete: access denied\n", fileID)
		return responses.NewEmptyResponse(http.StatusForbidden)
	}

	if err = deleteFile(persistedFile.Id); err != nil {
		log.Printf("File [%d] delete: %s\n", fileID, err.Error())
		return responses.NewEmptyResponse(http.StatusBadRequest)
	}

	log.Printf("File [%d] delete: success\n", fileID)
	return responses.NewEmptyResponse(http.StatusOK)
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
