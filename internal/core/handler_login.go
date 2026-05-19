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

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	cleanedUsername := strings.TrimSpace(req.Username)
	if cleanedUsername == "" || strings.TrimSpace(req.Password) == "" {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Username/password is required",
		)
		return
	}

	user, err := a.Queries.GetUserCredentialsByUsername(
		r.Context(),
		cleanedUsername,
	)
	if errors.Is(err, sql.ErrNoRows) {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusUnauthorized,
			"Invalid username/password",
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

	matches, err := argon2id.ComparePasswordAndHash(req.Password, user.Password)
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Something went wrong",
		)
		return
	}

	if !matches {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusUnauthorized,
			"Invalid username/password",
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

	refreshToken, err := a.insertRefreshToken(r.Context(), user.ID)
	if err != nil {
		a.Logger.Info("Failed to create refresh token", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to create refresh token",
		)
		return
	}

	a.setRefreshTokenCookie(w, *refreshToken)
	a.jsonResponse(w, http.StatusOK, accessToken)
}
