package web

import (
	"context"
	"github.com/lyubomirr/meme-generator-app/core/services"
	"net/http"
	"strings"
)

const Bearer = "Bearer "

func ValidateJwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, Bearer) {
			http.Error(w, "invalid jwt", http.StatusUnauthorized)
			return
		}

		jwt := authHeader[len(Bearer):]
		claims, err := tokenHandler.ValidateToken(jwt)
		if err != nil {
			http.Error(w, "invalid jwt", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), services.UserClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}