package files

import (
	"fmt"
	"go-heroku-server/src/responses"
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
		log.Printf("File create: %v\n", err)
		return responses.CreateResponse(http.StatusInternalServerError, nil)
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
	return responses.CreateResponse(http.StatusCreated, nil)
}

//Function provides requested file to the client
func readFile(userID uint, fileID uint) responses.IResponse {
	var persistedFile, err = getFile(fileID)
	if err != nil {
		log.Printf("File [%d] read: %v\n", fileID, err)
		return responses.CreateResponse(http.StatusNotFound, nil)
	}

	if err = persistedFile.validateAccess(userID); err != nil {
		log.Printf("File [%d] read: %v\n", fileID, err)
		return responses.CreateResponse(http.StatusForbidden, nil)
	}

	responseMap := map[string]string{
		"Access-Control-Expose-Headers": "Content-Disposition, Content-Length, X-Content-Transfer-Id",
		"Access-Control-Allow-Origin":   "*",
		"Content-Disposition":           "attachment; filename=" + persistedFile.FileName,
		"Content-Type":                  persistedFile.FileType,
	}

	log.Printf("File [%d] read: success\n", fileID)

	res := responses.CreateResponse(http.StatusOK, persistedFile.FileData)
	res.SetHeaders(responseMap)
	return res
}

func getFileList(userID uint) responses.IResponse {
	files := getAll(userID)
	for index := range files {
		fileName := files[index].FileName
		fileName = fileName[:strings.IndexByte(fileName, '.')]
		files[index].FileName = fileName
	}
	log.Printf("File listing: success\n")
	return responses.CreateResponse(http.StatusOK, files)
}

//Deletion of file base on userID
func removeFile(userID, fileID uint) responses.IResponse {
	persistedFile, err := getFile(fileID)
	if err != nil {
		log.Printf("File [%d] delete: %v\n", fileID, err)
		return responses.CreateResponse(http.StatusOK, nil)
	}

	if err = persistedFile.validateAccess(userID); err != nil {
		log.Printf("File [%d] delete: %v\n", fileID, err)
		return responses.CreateResponse(http.StatusForbidden, nil)
	}

	if err = deleteFile(persistedFile.Id); err != nil {
		log.Printf("File [%d] delete: %v\n", fileID, err)
		return responses.CreateResponse(http.StatusBadRequest, nil)
	}

	log.Printf("File [%d] delete: success\n", fileID)
	return responses.CreateResponse(http.StatusOK, nil)
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
