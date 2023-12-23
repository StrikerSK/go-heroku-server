package responses

import (
	"go-heroku-server/api/src/errors"
	"log"
	"net/http"
)

type ResponseFactory struct{}

func NewResponseFactory() ResponseFactory {
	return ResponseFactory{}
}

func (ResponseFactory) CreateResponse(input interface{}) (response IResponse) {
	switch input.(type) {
	case []byte:
		response = NewFileResponse(input.([]byte))
		break
	case errors.NotFoundError:
		response = NewEmptyResponse(http.StatusNotFound)
		break
	case errors.ConflictError:
		response = NewEmptyResponse(http.StatusConflict)
		break
	case errors.ForbiddenError:
		response = NewEmptyResponse(http.StatusForbidden)
		break
	case errors.UnauthorizedError:
		response = NewEmptyResponse(http.StatusUnauthorized)
		break
	case errors.ParseError, errors.DatabaseError:
		response = NewEmptyResponse(http.StatusBadRequest)
		break
	case errors.BadRequestError:
		response = NewBodyResponse(input)
		break
	case error:
		log.Println(input.(error))
		response = NewEmptyResponse(http.StatusInternalServerError)
		break
	case nil:
		response = NewEmptyResponse(http.StatusOK)
		break
	default:
		response = NewBodyResponse(input)
		break
	}
	return
}
