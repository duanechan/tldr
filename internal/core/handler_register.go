package core

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/duanechan/tldr/internal/auth"
	"github.com/duanechan/tldr/internal/database"
	"github.com/google/uuid"
)

type registerRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

const (
	minimumUsernameLength = 3
	minimumPasswordLength = 8
)

func (t *TLDR) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	cleanedUsername := strings.TrimSpace(req.Username)
	if cleanedUsername == "" || strings.TrimSpace(req.Password) == "" {
		errorResponse(w, http.StatusBadRequest, "Username/password is required")
		return
	}

	if len(cleanedUsername) < minimumUsernameLength {
		errorResponse(w, http.StatusBadRequest, "Username must be at least 3 characters long")
		return
	}

	if req.Password != req.ConfirmPassword {
		errorResponse(w, http.StatusBadRequest, "Passwords do not match")
		return
	}

	if len(req.Password) < minimumPasswordLength {
		errorResponse(w, http.StatusBadRequest, "Password must be at least 8 characters long")
		return
	}

	hashedPassword, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := t.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:       id,
		Username: req.Username,
		Password: hashedPassword,
	})
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, err := auth.CreateJWT(id, t.Config.JWTSecret, t.Config.JWTExpiry)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := t.insertRefreshToken(r.Context(), user)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	t.setRefreshTokenCookie(w, *refreshToken)
	jsonResponse(w, http.StatusCreated, authResponse{AccessToken: accessToken})
}
