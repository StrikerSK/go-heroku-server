package src

import (
	"encoding/json"
	"log"
	"net/http"
)

type RequestError struct {
	StatusCode   int    `json:"-"`
	ErrorMessage string `json:"error"`
}

func NewErrorResponse(statusCode int, error error) RequestError {
	return RequestError{
		StatusCode:   statusCode,
		ErrorMessage: error.Error(),
	}
}

func (re RequestError) WriteResponse(w http.ResponseWriter) {
	log.Println("Error: ", re.ErrorMessage)
	w.WriteHeader(re.StatusCode)
	payload, _ := json.Marshal(re)
	_, _ = w.Write(payload)
	return
}
