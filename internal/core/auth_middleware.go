package core

import (
	"context"
	"net/http"

	"github.com/duanechan/tldr/internal/auth"
)

type contextKey string

const claimsKey contextKey = "claims"

func (t *TLDR) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			errorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		claims, err := auth.ValidateJWT(token, t.Config.JWTSecret)
		if err != nil {
			errorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
