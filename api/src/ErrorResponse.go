package src

import (
	"encoding/json"
	"log"
	"net/http"
)

type RequestError struct {
	StatusCode   int               `json:"-"`
	ErrorMessage string            `json:"error"`
	Header       map[string]string `json:"-"`
}

func NewErrorResponse(statusCode int, error error) RequestError {
	return RequestError{
		StatusCode:   statusCode,
		ErrorMessage: error.Error(),
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func (re RequestError) WriteResponse(w http.ResponseWriter) {
	log.Println("Error: ", re.ErrorMessage)
	for key, value := range re.Header {
		w.Header().Add(key, value)
	}

	w.WriteHeader(re.StatusCode)
	payload, _ := json.Marshal(re)
	_, _ = w.Write(payload)
	return
}

func (re RequestError) AddHeader(newKey, keyValue string) {
	re.Header[newKey] = keyValue
	return
}
