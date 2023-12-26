package userHandlers

import (
	"context"
	"go-heroku-server/api/src/errors"
	"go-heroku-server/api/src/responses"
	userPorts "go-heroku-server/api/user/ports"
	"log"
	"net/http"
)

const (
	usernameContextKey       = "username"
	identificationContextKey = "id"
	rolesContextKey          = "roles"
)

type UserAuthMiddleware struct {
	tokenService    userPorts.ITokenService
	responseService responses.ResponseFactory
}

func NewUserAuthMiddleware(tokenService userPorts.ITokenService, responseService responses.ResponseFactory) UserAuthMiddleware {
	return UserAuthMiddleware{
		tokenService:    tokenService,
		responseService: responseService,
	}
}

func (h UserAuthMiddleware) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Println("Verify JWT Token: Cannot find header")
			h.responseService.CreateResponse(errors.NewUnauthorizedError("")).WriteResponse(w)
			return
		}

		userClaim, err := h.tokenService.ParseToken(token)
		if err != nil {
			h.responseService.CreateResponse(err).WriteResponse(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, usernameContextKey, userClaim.Username)
		ctx = context.WithValue(ctx, rolesContextKey, userClaim.Role)
		ctx = context.WithValue(ctx, identificationContextKey, userClaim.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h UserAuthMiddleware) GetUserIdentificationFromContext(context context.Context) (uint, error) {
	value, ok := context.Value(identificationContextKey).(uint)
	if !ok {
		return 0, errors.NewParseError("cannot resolve user from context")
	}

	return value, nil
}

func (h UserAuthMiddleware) GetUsernameFromContext(context context.Context) (string, error) {
	value, ok := context.Value(usernameContextKey).(string)
	if !ok {
		return "", errors.NewParseError("cannot resolve user from context")
	}

	return value, nil
}
