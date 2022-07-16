package userHandlers

import (
	"context"
	"errors"
	"go-heroku-server/api/src/responses"
	customAuth "go-heroku-server/api/user/auth"
	"log"
	"net/http"
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
		ctx := context.WithValue(r.Context(), userIdContextKey, userClaim.Username)
		ctx = context.WithValue(ctx, "roles", userClaim.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h UserAuthMiddleware) GetUserFromContext(context context.Context) (string, error) {
	value, ok := context.Value(userIdContextKey).(string)
	if !ok {
		return "", errors.New("cannot resolve user from context")
	}

	return value, nil
}
