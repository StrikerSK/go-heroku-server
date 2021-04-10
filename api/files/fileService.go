package files

import (
	"encoding/json"
	"go-heroku-server/api/user"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Function stores files received from the Front-End
func uploadFile(w http.ResponseWriter, r *http.Request) {

	var newFile File

	file, fileHeader, err := r.FormFile("file")

	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	fileBytes, _ := ioutil.ReadAll(file)
	defer file.Close()

	//Processing of received file metadata
	r.ParseForm()

	//r.FormValue("file") -> Receives metadata from sent file
	//Spracovanie JSONu na jednotlive polozky (Pozor! Je dolezite mat vytvorenu struct formu!!!)
	var token user.Token
	json.Unmarshal([]byte(r.FormValue("file")), &token)
	//fmt.Printf("hello: %s, testing: %s", metadata.Hello, metadata.Testing)

	newFile.FileName = fileHeader.Filename
	newFile.FileData = fileBytes
	newFile.FileSize = getFileSize(fileHeader.Size)
	newFile.FileType = http.DetectContentType(fileBytes)
	newFile.CreateDate = time.Now()

	claimId := r.Context().Value(user.UserIdContextKey)
	if claimId == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	newFile.UserID = claimId.(uint)
	createFile(newFile)
}

//Function provides requested file to the client
func readFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.Context().Value("file_ID").(int64)
	claimId := r.Context().Value(user.UserIdContextKey)
	if claimId == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var gotFile, err = getFile(fileID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf(err.Error() + " for id: " + strconv.FormatInt(fileID, 10))
		return
	}

	if gotFile.UserID != claimId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, Content-Length, X-Content-Transfer-Id")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Disposition", "attachment; filename="+gotFile.FileName)
	w.Header().Set("Content-Type", gotFile.FileType)

	_, _ = w.Write(gotFile.FileData)
}

func getFileList(w http.ResponseWriter, r *http.Request) {
	claimId := r.Context().Value(user.UserIdContextKey)
	if claimId == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	files := getAll(claimId)
	for index, file := range files {
		fileName := file.FileName
		fileName = fileName[:strings.IndexByte(fileName, '.')]
		files[index].FileName = fileName
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	_ = json.NewEncoder(w).Encode(files)
	log.Println("Retrieved list of file IDs and names")
}

//Function provides requested file to the client
func removeFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.Context().Value("file_ID").(int64)
	claimId := r.Context().Value(user.UserIdContextKey)
	if claimId == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	persistedFile, err := getFile(fileID)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		log.Printf(err.Error() + " for id: " + strconv.FormatInt(fileID, 10))
		return
	}

	if persistedFile.UserID != claimId {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	log.Printf("Deleted file with ID: %d", persistedFile.Id)
	_, _ = deleteFile(persistedFile.Id)
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

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(message))
}
