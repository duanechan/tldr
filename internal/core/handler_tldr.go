package core

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/duanechan/tldr/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (t *TLDR) UserGetTLDR(w http.ResponseWriter, r *http.Request) {
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

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse TLDR ID")
		return
	}

	tldr, err := t.Queries.GetTLDRByIDAndUser(r.Context(), database.GetTLDRByIDAndUserParams{
		UserID: userId,
		ID:     tldrId,
	})
	if errors.Is(err, sql.ErrNoRows) {
		t.errorResponse(w, r.Context(), http.StatusNotFound, fmt.Sprintf("TLDR with ID: %s not found", tldrId.String()))
		return
	}

	if err != nil {
		t.Logger.Error("Failed to get TLDR", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get TLDR with ID")
		return
	}

	t.jsonResponse(w, http.StatusOK, tldr)
}

func (t *TLDR) UserGetTLDRs(w http.ResponseWriter, r *http.Request) {
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

	tldrs, err := t.Queries.GetTLDRsByUser(r.Context(), userId)
	if errors.Is(err, sql.ErrNoRows) {
		t.jsonResponse(w, http.StatusOK, []database.Tldr{})
		return
	}

	if err != nil {
		t.Logger.Error("Failed to get TLDRs", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get TLDRs")
		return
	}

	t.jsonResponse(w, http.StatusOK, tldrs)
}

func (t *TLDR) AdminGetTLDR(w http.ResponseWriter, r *http.Request) {
	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse TLDR ID")
		return
	}

	tldr, err := t.Queries.GetTLDRById(r.Context(), tldrId)
	if errors.Is(err, sql.ErrNoRows) {
		t.errorResponse(w, r.Context(), http.StatusNotFound, fmt.Sprintf("TLDR with ID: %s not found", tldrId.String()))
		return
	}

	if err != nil {
		t.Logger.Error("Failed to get TLDR", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get TLDR")
		return
	}

	t.jsonResponse(w, http.StatusOK, tldr)
}

func (t *TLDR) AdminGetTLDRs(w http.ResponseWriter, r *http.Request) {
	tldrs, err := t.Queries.GetAllTLDRs(r.Context())
	if errors.Is(err, sql.ErrNoRows) {
		t.jsonResponse(w, http.StatusOK, []database.Tldr{})
		return
	}

	if err != nil {
		t.Logger.Error("Failed to get TLDRs", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get TLDRs")
		return
	}

	t.jsonResponse(w, http.StatusOK, tldrs)
}

func (t *TLDR) UserUpdateTLDR(w http.ResponseWriter, r *http.Request) {
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

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse TLDR ID")
		return
	}

	var updateRequest database.UpdateTLDRTitleParams
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Invalid request body")
		return
	}

	updateRequest.UserID = userId
	updateRequest.ID = tldrId

	tldr, err := t.Queries.UpdateTLDRTitle(r.Context(), updateRequest)
	if errors.Is(err, sql.ErrNoRows) {
		t.errorResponse(w, r.Context(), http.StatusNotFound, fmt.Sprintf("TLDR with ID: %s not found", tldrId.String()))
		return
	}

	if err != nil {
		t.Logger.Error("Failed to update TLDR", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to update TLDR")
		return
	}

	t.jsonResponse(w, http.StatusOK, tldr)
}

func (t *TLDR) AdminUpdateTLDR(w http.ResponseWriter, r *http.Request) {
	var updateRequest database.UpdateTLDRTitleParams
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Invalid request body")
		return
	}

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse TLDR ID")
		return
	}

	updateRequest.ID = tldrId

	tldr, err := t.Queries.UpdateTLDRTitle(r.Context(), updateRequest)
	if errors.Is(err, sql.ErrNoRows) {
		t.errorResponse(w, r.Context(), http.StatusNotFound, fmt.Sprintf("TLDR with ID: %s not found", tldrId))
		return
	}

	if err != nil {
		t.Logger.Error("Failed to update TLDR", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to update TLDR")
		return
	}

	t.jsonResponse(w, http.StatusOK, tldr)
}

func (t *TLDR) UserDeleteTLDR(w http.ResponseWriter, r *http.Request) {
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

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse TLDR ID")
		return
	}

	if err := t.Queries.DeleteTLDR(r.Context(), database.DeleteTLDRParams{
		UserID: userId,
		ID:     tldrId,
	}); err != nil {
		t.Logger.Error("Failed to delete TLDR", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to delete TLDR")
		return
	}

	t.jsonResponse(w, http.StatusNoContent, nil)
}

func (t *TLDR) AdminDeleteTLDR(w http.ResponseWriter, r *http.Request) {
	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse TLDR ID")
		return
	}

	if err = t.Queries.DeleteTLDRById(r.Context(), tldrId); err != nil {
		t.Logger.Error("Failed to delete TLDR", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to delete TLDR")
		return
	}

	t.jsonResponse(w, http.StatusNoContent, nil)
}
