package src

import (
	"fmt"
	"log"
	"net/http"
)

type RequestError struct {
	StatusCode int   `json:"-"`
	Err        error `json:"error"`
}

func (re RequestError) Error() string {
	return fmt.Sprintf("Status %d: err %v", re.StatusCode, re.Err)
}

func NewErrorResponse(statusCode int, err error) RequestError {
	return RequestError{
		StatusCode: statusCode,
		Err:        err,
	}
}

func (re RequestError) WriteResponse(w http.ResponseWriter) {
	log.Println("Error: ", re.Err)
	w.WriteHeader(re.StatusCode)
	//payload, _ := json.Marshal(re)
	//_, _ = w.Write(payload)
	return
}
