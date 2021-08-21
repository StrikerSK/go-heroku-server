package responses

import (
	"net/http"
)

type IResponse interface {
	WriteResponse(w http.ResponseWriter)
	AddHeader(newKey, keyValue string)
	SetHeaders(input map[string]string)
}

func CreateResponse(statusCode int, input interface{}) (res IResponse) {
	switch input.(type) {
	case []byte:
		res = NewFileResponse(input.([]byte))
		break
	case error:
		res = NewErrorResponse(statusCode, input.(error))
		break
	case nil:
		res = NewEmptyResponse(statusCode)
		break
	default:
		res = NewBodyResponse(input)
		break
	}
	return
}
