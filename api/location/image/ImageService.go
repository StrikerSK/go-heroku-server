package image

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func getLocationImage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uri, _ := strconv.ParseInt(vars["imageId"], 10, 64)

	var receivedImage = getImageFromDb(uri)

	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition, Content-Length, X-Content-Transfer-Id")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Disposition", "attachment; filename="+receivedImage.FileName)
	w.Header().Set("Content-Type", receivedImage.FileType)

	w.Write(receivedImage.FileData)
}
