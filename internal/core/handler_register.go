package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/duanechan/tldr/internal/auth"
	"github.com/duanechan/tldr/internal/database"
	"github.com/duanechan/tldr/internal/validate"
	"github.com/google/uuid"
	"modernc.org/sqlite"
)

const (
	minimumUsernameLength = 3
	minimumPasswordLength = 8

	sqliteUniqueConstraint = 2067
)

type registerRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

// Register handles new user registration requests by validating user
// credentials, setting refresh token cookie, and returning an access
// token.
func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	fieldErrors := []FieldError{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	username, errs := validate.String(
		req.Username,
		validate.Min(minimumUsernameLength),
		validate.NoWhitespace(),
	)
	for _, err := range errs {
		switch err {
		case validate.ErrMinLimit:
			fieldErrors = append(
				fieldErrors,
				FieldError{
					Field: "username",
					Message: fmt.Sprintf(
						"Username must be %d characters long",
						minimumUsernameLength,
					),
				},
			)
		case validate.ErrContainsWhitespace:
			fieldErrors = append(
				fieldErrors,
				FieldError{
					Field:   "username",
					Message: "Username must not contain whitespace.",
				},
			)
		default:
			fieldErrors = append(
				fieldErrors,
				FieldError{
					Field:   "username",
					Message: "Unknown validation error.",
				},
			)
		}
	}

	password, errs := validate.String(
		req.Password,
		validate.Min(minimumPasswordLength),
	)
	for _, err := range errs {
		switch err {
		case validate.ErrMinLimit:
			fieldErrors = append(
				fieldErrors,
				FieldError{
					Field: "password",
					Message: fmt.Sprintf(
						"Password must be %d characters long",
						minimumPasswordLength,
					),
				},
			)
		default:
			fieldErrors = append(
				fieldErrors,
				FieldError{
					Field:   "password",
					Message: "Unknown validation error.",
				},
			)
		}
	}

	if len(fieldErrors) > 0 {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to validate credentials",
			fieldErrors...)
		return
	}

	if password != req.ConfirmPassword {
		fieldErrors = append(
			fieldErrors,
			FieldError{
				Field:   "password",
				Message: "Passwords do not match",
			},
		)
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to validate credentials",
			fieldErrors...,
		)
		return
	}

	hashedPassword, err := argon2id.CreateHash(
		password,
		argon2id.DefaultParams,
	)
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Something went wrong",
		)
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Something went wrong",
		)
		return
	}

	user, err := a.Queries.CreateUser(r.Context(), database.CreateUserParams{
		ID:       id,
		Username: username,
		Password: hashedPassword,
	})
	if sqliteErr, ok := err.(*sqlite.Error); ok &&
		sqliteErr.Code() == sqliteUniqueConstraint {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusConflict,
			"Username already taken",
		)
		return
	}

	if err != nil {
		a.Logger.Error("Failed to create user", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to create user",
		)
		return
	}

	accessToken, err := auth.CreateJWT(
		id,
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
	a.jsonResponse(w, http.StatusCreated, accessToken)
}

// insertRefreshToken generates a refresh token, inserts it to the database,
// and returns it.
func (a *App) insertRefreshToken(
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

	refreshToken, err := a.Queries.CreateRefreshToken(
		ctx,
		database.CreateRefreshTokenParams{
			ID:        refreshTokenId,
			Token:     refreshTokenString,
			UserID:    id,
			ExpiresAt: time.Now().Add(a.Config.RefreshExpiry),
		},
	)
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

// setRefreshTokenCookie sets the "REFRESH_TOKEN" cookie.
func (a *App) setRefreshTokenCookie(
	w http.ResponseWriter,
	refreshToken database.RefreshToken,
) {
	var sameSite http.SameSite

	switch a.Config.Environment {
	case "prod":
		sameSite = http.SameSiteStrictMode
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
		Secure:   a.Config.Environment == "prod",
	})
}
