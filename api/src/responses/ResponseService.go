package responses

import (
	"go-heroku-server/api/src/errors"
	"log"
	"net/http"
)

type ResponseService struct{}

func NewResponseService() ResponseService {
	return ResponseService{}
}

func (ResponseService) CreateResponse(input interface{}) (response IResponse) {
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
