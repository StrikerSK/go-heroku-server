package responses

import (
	"net/http"
)

type FileResponse struct {
	StatusCode int               `json:"-"`
	Data       []byte            `json:"_"`
	Header     map[string]string `json:"-"`
}

func NewFileResponse(data []byte) *FileResponse {
	return &FileResponse{
		StatusCode: http.StatusOK,
		Data:       data,
	}
}

func (ir FileResponse) WriteResponse(w http.ResponseWriter) {
	for key, value := range ir.Header {
		w.Header().Add(key, value)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(ir.Data)
	return
}

func (ir *FileResponse) AddHeader(newKey, keyValue string) {
	ir.Header[newKey] = keyValue
	return
}

func (ir *FileResponse) SetHeaders(input map[string]string) {
	ir.Header = input
	return
}
