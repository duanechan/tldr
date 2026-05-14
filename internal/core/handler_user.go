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

func (t *TLDR) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsKey).(*jwt.RegisteredClaims)
	if !ok {
		t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid claims")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
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

func (t *TLDR) GetUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsKey).(*jwt.RegisteredClaims)
	if !ok {
		t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid claims")
		return
	}

	requestedUserId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		return
	}

	_, err = t.Queries.IsAdmin(r.Context(), userId)
	if err != nil && userId != requestedUserId {
		t.errorResponse(w, r.Context(), http.StatusForbidden, "You are not allowed to access or perform any actions to this resource")
		return
	}

	user, err := t.Queries.GetUserById(r.Context(), requestedUserId)
	if errors.Is(err, sql.ErrNoRows) {
		t.errorResponse(w, r.Context(), http.StatusNotFound, fmt.Sprintf("User with ID: %s not found", requestedUserId))
		return
	}

	if err != nil {
		t.Logger.Error("Failed to get user", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get user")
		return
	}

	t.jsonResponse(w, http.StatusOK, user)
}

func (t *TLDR) GetUsers(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsKey).(*jwt.RegisteredClaims)
	if !ok {
		t.errorResponse(w, r.Context(), http.StatusUnauthorized, "Invalid claims")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		return
	}

	_, err = t.Queries.IsAdmin(r.Context(), userId)
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusForbidden, "You are not allowed to access or perform any actions to this resource")
		return
	}

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
