package responses

import (
	"encoding/json"
	"net/http"
)

type ResponseImpl struct {
	Data   interface{}       `json:"data"`
	Header map[string]string `json:"-"`
}

func NewResponse(data interface{}) ResponseImpl {
	return ResponseImpl{
		Data: data,
		Header: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func (ri ResponseImpl) WriteResponse(w http.ResponseWriter) {
	for key, value := range ri.Header {
		w.Header().Add(key, value)
	}

	w.WriteHeader(http.StatusOK)
	payload, _ := json.Marshal(ri.Data)
	_, _ = w.Write(payload)
	return
}

func (ri ResponseImpl) AddHeader(newKey, keyValue string) {
	ri.Header[newKey] = keyValue
	return
}
