package image

import (
	"go-heroku-server/src/responses"
	"log"
	"net/http"
)

func sReadLocationImage(imageID int64) responses.IResponse {
	persistedImage, err := readImage(imageID)
	if err != nil {
		log.Printf("Location image [%d] read: %s\n", imageID, err)
		return responses.CreateResponse(http.StatusBadRequest, nil)
	}

	responseMap := map[string]string{
		"Access-Control-Expose-Headers": "Content-Disposition, Content-Length, X-Content-Transfer-Id",
		"Access-Control-Allow-Origin":   "*",
		"Content-Disposition":           "attachment; filename=" + persistedImage.FileName,
		"Content-Type":                  persistedImage.FileType,
	}

	log.Printf("Location image [%d] read: success\n", imageID)
	res := responses.CreateResponse(http.StatusOK, persistedImage.FileData)
	res.SetHeaders(responseMap)
	return res
}
