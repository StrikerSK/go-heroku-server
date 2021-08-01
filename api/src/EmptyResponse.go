package src

import (
	"net/http"
)

type EmptyResponse int

func NewEmptyResponse(statusCode int) EmptyResponse {
	return EmptyResponse(statusCode)
}

func (ri EmptyResponse) WriteResponse(w http.ResponseWriter) {
	w.WriteHeader(int(ri))
	return
}
