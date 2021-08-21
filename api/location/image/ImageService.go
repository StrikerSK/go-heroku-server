package image

import (
	"go-heroku-server/api/src/responses"
	"log"
	"net/http"
)

func sReadLocationImage(imageID int64) responses.IResponse {
	persistedImage, err := readImage(imageID)
	if err != nil {
		log.Printf("Location image [%d] read: %s\n", imageID, err)
		return responses.NewEmptyResponse(http.StatusBadRequest)
	}

	responseMap := map[string]string{
		"Access-Control-Expose-Headers": "Content-Disposition, Content-Length, X-Content-Transfer-Id",
		"Access-Control-Allow-Origin":   "*",
		"Content-Disposition":           "attachment; filename=" + persistedImage.FileName,
		"Content-Type":                  persistedImage.FileType,
	}

	log.Printf("Location image [%d] read: success\n", imageID)
	return responses.NewFileResponse(persistedImage.FileData, responseMap)
}
