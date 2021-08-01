package src

import (
	"encoding/json"
	"net/http"
)

type ResponseImpl struct {
	Data interface{} `json:"data"`
}

func NewResponse(data interface{}) ResponseImpl {
	return ResponseImpl{
		Data: data,
	}
}

func (ri ResponseImpl) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	payload, _ := json.Marshal(ri.Data)
	_, _ = w.Write(payload)
	return
}
