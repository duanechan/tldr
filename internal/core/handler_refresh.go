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
		t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid or missing cookie")
		return
	}

	if err := cookie.Valid(); err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Invalid or missing cookie")
		return
	}

	token := cookie.Value

	user, err := t.Queries.GetUserByRefreshToken(r.Context(), token)
	if errors.Is(err, sql.ErrNoRows) {
		t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	if err != nil {
		t.Logger.Error("Failed to get user", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get user")
		return
	}

	accessToken, err := auth.CreateJWT(user.ID, t.Config.JWTSecret, t.Config.JWTExpiry)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to create access token")
		return
	}

	t.jsonResponse(w, http.StatusOK, authResponse{AccessToken: accessToken})
}
