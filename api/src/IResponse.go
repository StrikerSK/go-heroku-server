package src

import "net/http"

type IResponse interface {
	WriteResponse(w http.ResponseWriter)
}
