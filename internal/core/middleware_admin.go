package core

import (
	"net/http"

	"github.com/duanechan/tldr/internal/auth"
)

func (a *App) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, err := auth.GetUserID(r.Context())
		if err != nil {
			a.errorResponse(
				w,
				r.Context(),
				http.StatusUnauthorized,
				"Invalid claims",
			)
			return
		}

		if _, err = a.Queries.IsAdmin(r.Context(), userId); err != nil {
			a.errorResponse(
				w,
				r.Context(),
				http.StatusForbidden,
				"You are not allowed to access or perform any actions to this resource",
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
