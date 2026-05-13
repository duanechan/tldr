package core

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/duanechan/tldr/internal/auth"
	"github.com/duanechan/tldr/internal/database"
	"github.com/google/uuid"
	"modernc.org/sqlite"
)

const (
	minimumUsernameLength = 3
	minimumPasswordLength = 8
)

func (t *TLDR) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cleanedUsername := strings.TrimSpace(req.Username)
	if cleanedUsername == "" || strings.TrimSpace(req.Password) == "" {
		t.errorResponse(w, http.StatusBadRequest, "Username/password is required")
		return
	}

	if len(cleanedUsername) < minimumUsernameLength {
		t.errorResponse(w, http.StatusBadRequest, "Username must be at least 3 characters long")
		return
	}

	if req.Password != req.ConfirmPassword {
		t.errorResponse(w, http.StatusBadRequest, "Passwords do not match")
		return
	}

	if len(req.Password) < minimumPasswordLength {
		t.errorResponse(w, http.StatusBadRequest, "Password must be at least 8 characters long")
		return
	}

	hashedPassword, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		t.errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		t.errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user, err := t.Queries.CreateUser(r.Context(), database.CreateUserParams{
		ID:       id,
		Username: cleanedUsername,
		Password: hashedPassword,
	})
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok {
			if sqliteErr.Code() == 2067 {
				t.errorResponse(w, http.StatusConflict, "Username already taken")
				return
			}
		}
		t.Logger.Error("Failed to create user", "error", err.Error())
		t.errorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	accessToken, err := auth.CreateJWT(id, t.Config.JWTSecret, t.Config.JWTExpiry)
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
	t.jsonResponse(w, http.StatusCreated, authResponse{AccessToken: accessToken})
}
