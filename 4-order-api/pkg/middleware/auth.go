package middleware

import (
	"4-order-api/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type key string

const (
	UserIdContextKey key    = "user_id"
	authPrefix       string = "Bearer "
)

func writeUnauthed(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func IsAuthed(next http.Handler, JWT *jwt.JWT) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, authPrefix) {
			writeUnauthed(w)
			return
		}

		token := strings.TrimPrefix(authHeader, authPrefix)
		isValid, data := JWT.Parse(token)
		if !isValid {
			writeUnauthed(w)
			return
		}

		ctx := context.WithValue(r.Context(), UserIdContextKey, data.UserID)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

func GetUserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIdContextKey).(uint)
	return userID, ok
}
