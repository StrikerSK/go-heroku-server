package responses

import (
	"net/http"
)

type EmptyResponse struct {
	StatusCode int               `json:"-"`
	Header     map[string]string `json:"-"`
}

func NewEmptyResponse(statusCode int) EmptyResponse {
	return EmptyResponse{StatusCode: statusCode}
}

func (er EmptyResponse) WriteResponse(w http.ResponseWriter) {
	for key, value := range er.Header {
		w.Header().Add(key, value)
	}

	w.WriteHeader(er.StatusCode)
	return
}

func (er EmptyResponse) AddHeader(newKey, keyValue string) {
	er.Header[newKey] = keyValue
	return
}
