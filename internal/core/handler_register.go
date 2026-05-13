package core

import (
	"encoding/json"
	"fmt"
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
	fieldErrors := []FieldError{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Invalid request body")
		return
	}

	cleanedUsername := strings.TrimSpace(req.Username)
	if cleanedUsername == "" {
		fieldErrors = append(fieldErrors, FieldError{Field: "username", Message: "Username is required"})
	} else if len(cleanedUsername) < minimumUsernameLength {
		fieldErrors = append(fieldErrors, FieldError{Field: "username", Message: fmt.Sprintf("Username must be %d characters long", minimumUsernameLength)})
	}

	if strings.TrimSpace(req.Password) == "" {
		fieldErrors = append(fieldErrors, FieldError{Field: "password", Message: "Password is required"})
	} else if len(req.Password) < minimumPasswordLength {
		fieldErrors = append(fieldErrors, FieldError{Field: "password", Message: fmt.Sprintf("Password must be %d characters long", minimumPasswordLength)})
	} else if req.Password != req.ConfirmPassword {
		fieldErrors = append(fieldErrors, FieldError{Field: "password", Message: "Passwords do not match"})
	}

	if len(fieldErrors) > 0 {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to validate credentials", fieldErrors...)
		return
	}

	hashedPassword, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
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
				t.errorResponse(w, r.Context(), http.StatusConflict, "Username already taken")
				return
			}
		}
		t.Logger.Error("Failed to create user", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to create user")
		return
	}

	accessToken, err := auth.CreateJWT(id, t.Config.JWTSecret, t.Config.JWTExpiry)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to create access token")
		return
	}

	refreshToken, err := t.insertRefreshToken(r.Context(), user)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to create refresh token")
		return
	}
	t.setRefreshTokenCookie(w, *refreshToken)
	t.jsonResponse(w, http.StatusCreated, authResponse{AccessToken: accessToken})
}
