package core

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/duanechan/tldr/internal/auth"
)

func (t *TLDR) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("REFRESH_TOKEN")
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid or missing cookie")
		return
	}

	if err := cookie.Valid(); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid or missing cookie")
		return
	}

	token := cookie.Value

	user, err := t.Queries.GetUserByRefreshToken(r.Context(), token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}
		t.Logger.Error("Failed to get user", "error", err.Error())
		errorResponse(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	accessToken, err := auth.CreateJWT(user.ID, t.Config.JWTSecret, t.Config.JWTExpiry)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to create access token")
		return
	}

	jsonResponse(w, http.StatusOK, authResponse{AccessToken: accessToken})
}
