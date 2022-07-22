package userHandlers

import (
	"context"
	"go-heroku-server/api/src/errors"
	"go-heroku-server/api/src/responses"
	customAuth "go-heroku-server/api/user/service"
	"log"
	"net/http"
)

const (
	usernameContextKey       = "username"
	identificationContextKey = "id"
)

type UserAuthMiddleware struct {
	tokenService customAuth.TokenService
}

func NewUserAuthMiddleware(tokenService customAuth.TokenService) UserAuthMiddleware {
	return UserAuthMiddleware{
		tokenService: tokenService,
	}
}

func (h UserAuthMiddleware) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			log.Println("Verify JWT Token: Cannot find header")
			responses.CreateResponse(http.StatusUnauthorized, nil).WriteResponse(w)
			return
		}
		userClaim, err := h.tokenService.ParseToken(token)
		if err != nil {
			responses.CreateResponse(http.StatusUnauthorized, nil).WriteResponse(w)
			return
		}
		ctx := context.WithValue(r.Context(), usernameContextKey, userClaim.Username)
		ctx = context.WithValue(ctx, "roles", userClaim.Role)
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
