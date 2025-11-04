package middleware

import (
	"context"
	"linkr/internal/auth"
	"log"
	"net/http"
	"strings"
)

type ctxKey string

const UserCtxKey ctxKey = "user"

func AuthBearerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		log.Printf("Authorization %s", header)
		if header == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]
		claims, err := auth.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserCtxKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserFromContext(r *http.Request) *auth.Claims {
	claims, ok := r.Context().Value(UserCtxKey).(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}
