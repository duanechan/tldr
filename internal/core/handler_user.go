package core

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/duanechan/tldr/internal/database"
	"github.com/duanechan/tldr/internal/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (t *TLDR) UserGetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(types.ClaimsKey).(*jwt.RegisteredClaims)
	if !ok {
		t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid claims")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse user ID")
		return
	}

	t.getUser(w, r, userId)
}

func (t *TLDR) UserUpdateUsername(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(types.ClaimsKey).(*jwt.RegisteredClaims)
	if !ok {
		t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid claims")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to parse user ID")
		return
	}

	t.updateUsername(w, r, userId)
}

func (t *TLDR) UserUpdatePassword(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(types.ClaimsKey).(*jwt.RegisteredClaims)
	if !ok {
		t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid claims")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to parse user ID")
		return
	}

	t.updatePassword(w, r, userId)
}

func (t *TLDR) AdminGetUser(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse user ID")
		return
	}

	t.getUser(w, r, userId)
}

func (t *TLDR) AdminGetUsers(w http.ResponseWriter, r *http.Request) {
	cursor, limit, fieldErrors := extractQueryParams(r.URL.Query())
	if fieldErrors != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse query params", fieldErrors...)
		return
	}

	users, err := t.Queries.GetUsers(r.Context(), database.GetUsersParams{
		CreatedAt: time.Time(cursor),
		Limit:     int64(limit),
	})
	if errors.Is(err, sql.ErrNoRows) || users == nil {
		t.jsonResponse(w, http.StatusOK, []database.User{})
		return
	}

	if err != nil {
		t.Logger.Error("Failed to get users", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get users")
		return
	}

	if int(limit) > len(users) {
		t.jsonResponse(w, http.StatusOK, types.Page[database.GetUsersRow]{Results: users})
		return
	}

	next := users[limit-1]
	t.jsonResponse(w, http.StatusOK, types.Page[database.GetUsersRow]{
		Results: users[:limit-1],
		Next:    (*types.PageCursor)(&next.CreatedAt),
	})
}

func (t *TLDR) AdminUpdateUsername(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse user ID")
		return
	}

	t.updateUsername(w, r, userId)
}

func (t *TLDR) AdminUpdatePassword(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse user ID")
		return
	}

	t.updatePassword(w, r, userId)
}

func (t *TLDR) AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse user ID")
		return
	}

	if err = t.Queries.DeleteUser(r.Context(), userId); err != nil {
		t.Logger.Error("Failed to delete user", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to delete user")
		return
	}

	t.jsonResponse(w, http.StatusNoContent, nil)
}

func (t *TLDR) getUser(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	user, err := t.Queries.GetUserById(r.Context(), userId)
	if errors.Is(err, sql.ErrNoRows) {
		t.errorResponse(w, r.Context(), http.StatusNotFound, "User not found")
		return
	}

	if err != nil {
		t.Logger.Error("Failed to get user", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get user")
		return
	}

	t.jsonResponse(w, http.StatusOK, user)
}

func (t *TLDR) updateUsername(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	var updateRequest database.UpdateUsernameParams
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Invalid request body")
		return
	}

	updateRequest.ID = userId

	var fieldError *types.FieldError
	cleanedUsername := strings.TrimSpace(updateRequest.Username)
	if cleanedUsername == "" {
		fieldError = &types.FieldError{Field: "username", Message: "Username is required"}
	} else if len(cleanedUsername) < types.MinimumUsernameLength {
		fieldError = &types.FieldError{Field: "username", Message: fmt.Sprintf("Username must be %d characters long", types.MinimumUsernameLength)}
	}

	if fieldError != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to update username", *fieldError)
		return
	}

	updateRequest.Username = cleanedUsername

	user, err := t.Queries.UpdateUsername(r.Context(), updateRequest)
	if errors.Is(err, sql.ErrNoRows) {
		t.errorResponse(w, r.Context(), http.StatusNotFound, "Failed to update user")
		return
	}

	if err != nil {
		t.Logger.Error("Failed to update user", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to update user")
		return
	}

	t.jsonResponse(w, http.StatusOK, user)
}

func (t *TLDR) updatePassword(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
	var updateRequest database.UpdatePasswordParams
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Invalid request body")
		return
	}

	updateRequest.ID = userId

	var fieldError *types.FieldError
	if strings.TrimSpace(updateRequest.Password) == "" {
		fieldError = &types.FieldError{Field: "password", Message: "Password is required"}
	} else if len(updateRequest.Password) < types.MinimumPasswordLength {
		fieldError = &types.FieldError{Field: "password", Message: fmt.Sprintf("Password must be %d characters long", types.MinimumPasswordLength)}
	}

	if fieldError != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to update password", *fieldError)
		return
	}

	hashedPassword, err := argon2id.CreateHash(updateRequest.Password, argon2id.DefaultParams)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		return
	}

	updateRequest.Password = hashedPassword

	user, err := t.Queries.UpdatePassword(r.Context(), updateRequest)
	if errors.Is(err, sql.ErrNoRows) {
		t.errorResponse(w, r.Context(), http.StatusNotFound, "User not found")
		return
	}

	if err != nil {
		t.Logger.Error("Failed to update password", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to update password")
		return
	}

	t.jsonResponse(w, http.StatusOK, user)
}
