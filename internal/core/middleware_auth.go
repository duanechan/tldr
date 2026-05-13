package core

import (
	"context"
	"net/http"

	"github.com/duanechan/tldr/internal/auth"
)

func (t *TLDR) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Failed to get bearer token")
			return
		}

		claims, err := auth.ValidateJWT(token, t.Config.JWTSecret)
		if err != nil {
			t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
