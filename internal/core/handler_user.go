package core

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/duanechan/tldr/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (t *TLDR) UserGetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsKey).(*jwt.RegisteredClaims)
	if !ok {
		t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid claims")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse user ID")
		return
	}

	user, err := t.Queries.GetUserById(r.Context(), userId)
	if errors.Is(err, sql.ErrNoRows) {
		t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid session")
		return
	}

	if err != nil {
		t.Logger.Error("Failed to get user", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get user")
		return
	}

	t.jsonResponse(w, http.StatusOK, user)
}

func (t *TLDR) UserUpdateUsername(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsKey).(*jwt.RegisteredClaims)
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
	claims, ok := r.Context().Value(claimsKey).(*jwt.RegisteredClaims)
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

	user, err := t.Queries.GetUserById(r.Context(), userId)
	if errors.Is(err, sql.ErrNoRows) {
		t.errorResponse(w, r.Context(), http.StatusNotFound, fmt.Sprintf("User with ID: %s not found", userId))
		return
	}

	if err != nil {
		t.Logger.Error("Failed to get user", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get user")
		return
	}

	t.jsonResponse(w, http.StatusOK, user)
}

func (t *TLDR) AdminGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := t.Queries.GetUsers(r.Context())
	if errors.Is(err, sql.ErrNoRows) {
		t.jsonResponse(w, http.StatusOK, []database.User{})
		return
	}

	if err != nil {
		t.Logger.Error("Failed to get users", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get users")
		return
	}

	t.jsonResponse(w, http.StatusOK, users)
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
