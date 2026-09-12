package middleware

import (
	"context"
	"net/http"
	"strings"

	"blog-platform/pkg/auth"
	"blog-platform/pkg/httpjson"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			httpjson.WriteError(w, http.StatusUnauthorized, "missing or invalid token")
			return
		}

		claims, err := auth.ValidateToken(strings.TrimPrefix(tokenStr, "Bearer "))
		if err != nil {
			httpjson.WriteError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
