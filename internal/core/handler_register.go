package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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
		t.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	cleanedUsername := strings.TrimSpace(req.Username)
	if cleanedUsername == "" {
		fieldErrors = append(
			fieldErrors,
			FieldError{Field: "username", Message: "Username is required"},
		)
	} else if len(cleanedUsername) < minimumUsernameLength {
		fieldErrors = append(fieldErrors, FieldError{Field: "username", Message: fmt.Sprintf("Username must be %d characters long", minimumUsernameLength)})
	}

	if strings.TrimSpace(req.Password) == "" {
		fieldErrors = append(
			fieldErrors,
			FieldError{Field: "password", Message: "Password is required"},
		)
	} else if len(req.Password) < minimumPasswordLength {
		fieldErrors = append(fieldErrors, FieldError{Field: "password", Message: fmt.Sprintf("Password must be %d characters long", minimumPasswordLength)})
	} else if req.Password != req.ConfirmPassword {
		fieldErrors = append(fieldErrors, FieldError{Field: "password", Message: "Passwords do not match"})
	}

	if len(fieldErrors) > 0 {
		t.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to validate credentials",
			fieldErrors...)
		return
	}

	hashedPassword, err := argon2id.CreateHash(
		req.Password,
		argon2id.DefaultParams,
	)
	if err != nil {
		t.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Something went wrong",
		)
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		t.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Something went wrong",
		)
		return
	}

	user, err := t.Queries.CreateUser(r.Context(), database.CreateUserParams{
		ID:       id,
		Username: cleanedUsername,
		Password: hashedPassword,
	})
	if sqliteErr, ok := err.(*sqlite.Error); ok &&
		sqliteErr.Code() == sqliteUniqueConstraint {
		t.errorResponse(
			w,
			r.Context(),
			http.StatusConflict,
			"Username already taken",
		)
		return
	}

	if err != nil {
		t.Logger.Error("Failed to create user", "error", err.Error())
		t.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to create user",
		)
		return
	}

	accessToken, err := auth.CreateJWT(
		id,
		t.Config.JWTSecret,
		t.Config.JWTExpiry,
	)
	if err != nil {
		t.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to create access token",
		)
		return
	}

	refreshToken, err := t.insertRefreshToken(r.Context(), user.ID)
	if err != nil {
		t.Logger.Info("Failed to create refresh token", "error", err.Error())
		t.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to create refresh token",
		)
		return
	}

	t.setRefreshTokenCookie(w, *refreshToken)
	t.jsonResponse(
		w,
		http.StatusCreated,
		authResponse{AccessToken: accessToken},
	)
}

func (t *TLDR) insertRefreshToken(
	ctx context.Context,
	id uuid.UUID,
) (*database.RefreshToken, error) {
	refreshTokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshToken, err := t.Queries.CreateRefreshToken(
		ctx,
		database.CreateRefreshTokenParams{
			ID:        refreshTokenId,
			Token:     refreshTokenString,
			UserID:    id,
			ExpiresAt: time.Now().Add(t.Config.RefreshExpiry),
		},
	)
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (t *TLDR) setRefreshTokenCookie(
	w http.ResponseWriter,
	refreshToken database.RefreshToken,
) {
	var sameSite http.SameSite

	switch t.Config.Environment {
	case "prod":
		sameSite = http.SameSiteStrictMode
	case "dev":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteLaxMode
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "REFRESH_TOKEN",
		Value:    refreshToken.Token,
		Path:     "/",
		Expires:  refreshToken.ExpiresAt,
		HttpOnly: true,
		SameSite: sameSite,
		Secure:   t.Config.Environment == "prod",
	})
}
