package core

import (
	"net/http"

	"github.com/duanechan/tldr/internal/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (t *TLDR) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(types.ClaimsKey).(*jwt.RegisteredClaims)
		if !ok {
			t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid claims")
			return
		}

		userId, err := uuid.Parse(claims.Subject)
		if err != nil {
			t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
			return
		}

		if _, err := t.Queries.IsAdmin(r.Context(), userId); err != nil {
			t.errorResponse(w, r.Context(), http.StatusForbidden, "You are not allowed to access or perform any actions to this resource")
			return
		}

		next.ServeHTTP(w, r)
	})
}
