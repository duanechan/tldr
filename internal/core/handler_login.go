package core

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/duanechan/tldr/internal/auth"
)

func (t *TLDR) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cleanedUsername := strings.TrimSpace(req.Username)
	if cleanedUsername == "" || strings.TrimSpace(req.Password) == "" {
		t.errorResponse(w, http.StatusBadRequest, "Username/password is required")
		return
	}

	user, err := t.Queries.GetUserByName(r.Context(), cleanedUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			t.errorResponse(w, http.StatusUnauthorized, "Invalid username/password")
			return
		}
		t.Logger.Error("Failed to get user", "error", err.Error())
		t.errorResponse(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	matches, err := argon2id.ComparePasswordAndHash(req.Password, user.Password)
	if err != nil {
		t.errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if !matches {
		t.errorResponse(w, http.StatusUnauthorized, "Invalid username/password")
		return
	}

	accessToken, err := auth.CreateJWT(user.ID, t.Config.JWTSecret, t.Config.JWTExpiry)
	if err != nil {
		t.errorResponse(w, http.StatusInternalServerError, "Failed to create access token")
		return
	}

	refreshToken, err := t.insertRefreshToken(r.Context(), user)
	if err != nil {
		t.errorResponse(w, http.StatusInternalServerError, "Failed to create refresh token")
		return
	}

	t.setRefreshTokenCookie(w, *refreshToken)
	t.jsonResponse(w, http.StatusOK, authResponse{AccessToken: accessToken})
}
