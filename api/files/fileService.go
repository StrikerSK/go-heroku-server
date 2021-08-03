package files

import (
	"errors"
	"go-heroku-server/api/src"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Function stores files received from the Front-End
func uploadFile(file multipart.File, fileHeader *multipart.FileHeader, userID uint) {
	fileBytes, _ := ioutil.ReadAll(file)
	defer file.Close()

	customFile := File{
		UserID:     userID,
		FileName:   fileHeader.Filename,
		FileType:   http.DetectContentType(fileBytes),
		FileData:   fileBytes,
		FileSize:   getFileSize(fileHeader.Size),
		CreateDate: time.Now(),
	}

	createFile(customFile)
}

//Function provides requested file to the client
func readFile(userID uint, fileID uint) (*File, src.IResponse) {
	var gotFile, err = getFile(fileID)
	if err != nil {
		log.Printf(err.Error() + " for id: " + strconv.Itoa(int(fileID)))
		return nil, src.NewErrorResponse(http.StatusNotFound, err)
	}

	if gotFile.UserID != userID {
		log.Printf("access denied for file id: " + strconv.Itoa(int(fileID)))
		return nil, src.NewErrorResponse(http.StatusForbidden, err)
	}

	return &gotFile, nil
}

func getFileList(userID uint) src.ResponseImpl {
	files := getAll(userID)
	for index, file := range files {
		fileName := file.FileName
		fileName = fileName[:strings.IndexByte(fileName, '.')]
		files[index].FileName = fileName
	}
	return src.NewResponse(files)
}

//Function provides requested file to the client
func removeFile(userID, fileID uint) src.IResponse {
	persistedFile, err := getFile(fileID)
	if err != nil {
		log.Printf("File [%d] not found", fileID)
		return src.NewEmptyResponse(http.StatusOK)
	}

	if persistedFile.UserID != userID {
		err = errors.New("access denied")
		return src.NewErrorResponse(http.StatusForbidden, err)
	}

	if err = deleteFile(persistedFile.Id); err != nil {
		return src.NewErrorResponse(http.StatusBadRequest, err)
	}

	log.Printf("Deleted file with ID: %d", persistedFile.Id)
	return src.NewEmptyResponse(http.StatusOK)
}

func getFileSize(fileSize int64) (outputSize string) {

	switch {
	case fileSize < 1024:
		outputSize = strconv.FormatInt(fileSize, 10) + " B"
		break
	case fileSize < 1048576:
		fileSize = fileSize / 1024
		outputSize = strconv.FormatInt(fileSize, 10) + " kB"
		break
	default:
		fileSize = fileSize / 1048576
		outputSize = strconv.FormatInt(fileSize, 10) + " MB"
		break
	}

	return
}
