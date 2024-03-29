package responses

import (
	"net/http"
)

type IResponse interface {
	WriteResponse(w http.ResponseWriter)
	AddHeader(newKey, keyValue string)
	SetHeaders(input map[string]string)
}
