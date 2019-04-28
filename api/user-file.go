package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-heroku-server/api/types"
	"go-heroku-server/api/utils"
	"go-heroku-server/config"
)

type File struct {
	Id         int64     `json:"id"`
	UserID     uint      `json:"-"`
	FileName   string    `json:"filename"`
	FileType   string    `json:"-"`
	FileData   []byte    `json:"-"`
	FileSize   string    `json:"fileSize"`
	CreateDate time.Time `json:"createDate"`
}

type Metadata struct {
	Hello   string
	Testing string
}

//Function stores files received from the Front-End
func UploadFile(w http.ResponseWriter, r *http.Request) {

	var newFile File

	file, fileHeader, err := r.FormFile("file")

	if err != nil {
		renderError(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	fileBytes, _ := ioutil.ReadAll(file)

	defer file.Close()

	fileName := fileHeader.Filename

	//Tato cast ma na starosti spracovanie zaslanych metadat
	r.ParseForm()

	//Tato cast ziskava metadata zaslane so suborom
	fileData := r.FormValue("file")

	//Spracovanie JSONu na jednotlive polozky (Pozor! Je dolezite mat vytvorenu struct formu!!!)
	var token types.Token
	json.Unmarshal([]byte(fileData), &token)
	//fmt.Printf("hello: %s, testing: %s", metadata.Hello, metadata.Testing)

	newFile.FileName = fileName
	newFile.FileData = fileBytes
	newFile.FileSize = getFileSize(fileHeader.Size)
	newFile.FileType = http.DetectContentType(fileBytes)
	newFile.CreateDate = time.Now()
	newFile.UserID, _ = utils.GetIdFromToken(token.Token)
	//newFile.UserID, _ = utils.GetIdFromToken(fileData)

	processFile(newFile)
}

//Function provides requested file to the client
func ServeFile(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uri, _ := strconv.ParseInt(vars["id"], 10, 64)

	var gotFile = getFileFromDb(uri)

	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, Content-Length, X-Content-Transfer-Id")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Disposition", "attachment; filename="+gotFile.FileName)
	w.Header().Set("Content-Type", gotFile.FileType)

	w.Write(gotFile.FileData)
}

func processFile(file File) {

	db, err := config.CreateDatabase()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.NewRecord(file)
	db.Create(&file)

	log.Println("Inserted file: " + file.FileName + ".")

}

func getFileFromDb(fileId int64) File {

	db, err := config.CreateDatabase()
	if err != nil {
		panic(err)
	}

	var fileFile File

	db.Where("id = ?", fileId).Find(&fileFile)

	return fileFile

}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func GetFileList(w http.ResponseWriter, r *http.Request) {

	var files []File
	receivedToken := (r.Header.Get("Authorization"))

	if receivedToken != "null" {
		userId, _ := utils.GetIdFromToken(receivedToken)

		db, err := config.CreateDatabase()

		if err != nil {
			panic(err)
		}

		defer db.Close()

		db.Where("user_id = ?", userId).Find(&files)

		for index, file := range files {

			fileName := file.FileName
			fileName = fileName[:strings.IndexByte(fileName, '.')]

			files[index].FileName = fileName

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		json.NewEncoder(w).Encode(files)
		log.Println("Retrieved list of file IDs and names")
	} else {
		w.WriteHeader(400)
	}

}

func getFileSize(fileSize int64) string {

	officialFileSize := ""

	switch {
	case fileSize < 1024:
		officialFileSize = strconv.FormatInt(fileSize, 10) + " B"
		break
	case fileSize < 1048576:
		fileSize = fileSize / 1024
		officialFileSize = strconv.FormatInt(fileSize, 10) + " kB"
		break
	default:
		fileSize = fileSize / 1048576
		officialFileSize = strconv.FormatInt(fileSize, 10) + " MB"
		break
	}

	return officialFileSize
}
