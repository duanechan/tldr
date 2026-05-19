package core

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/duanechan/tldr/internal/auth"
)

func (a *App) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("REFRESH_TOKEN")
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusUnauthorized,
			"Invalid or missing cookie",
		)
		return
	}

	if err := cookie.Valid(); err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid or missing cookie",
		)
		return
	}

	token := cookie.Value

	user, err := a.Queries.GetUserByRefreshToken(r.Context(), token)
	if errors.Is(err, sql.ErrNoRows) {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusUnauthorized,
			"Invalid or expired token",
		)
		return
	}

	if err != nil {
		a.Logger.Error("Failed to get user", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to get user",
		)
		return
	}

	accessToken, err := auth.CreateJWT(
		user.ID,
		a.Config.JWTSecret,
		a.Config.JWTExpiry,
	)
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to create access token",
		)
		return
	}

	a.jsonResponse(w, http.StatusOK, accessToken)
}
