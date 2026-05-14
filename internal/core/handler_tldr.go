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

func (t *TLDR) GetTLDR(w http.ResponseWriter, r *http.Request) {
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

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		return
	}

	tldr, err := t.Queries.GetTLDRFromUser(r.Context(), database.GetTLDRFromUserParams{
		UserID: userId,
		ID:     tldrId,
	})
	if errors.Is(err, sql.ErrNoRows) {
		t.errorResponse(w, r.Context(), http.StatusNotFound, fmt.Sprintf("TLDR with ID: %s not found", tldrId.String()))
		return
	}

	if err != nil {
		t.Logger.Error("Failed to get TLDR", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get TLDR with ID: "+tldrId.String())
		return
	}

	t.jsonResponse(w, http.StatusOK, tldr)
}

func (t *TLDR) GetTLDRs(w http.ResponseWriter, r *http.Request) {
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

	tldrs, err := t.Queries.GetTLDRsFromUser(r.Context(), userId)
	if err != nil {
		t.Logger.Error("Failed to get TLDRs", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get TLDRs")
		return
	}

	t.jsonResponse(w, http.StatusOK, tldrs)
}

func (t *TLDR) UpdateTLDR(w http.ResponseWriter, r *http.Request) {
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

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
		return
	}

	updateRequest := database.UpdateTLDRTitleParams{
		UserID: userId,
		ID:     tldrId,
	}
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Invalid request body")
		return
	}

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

func (t *TLDR) DeleteTLDR(w http.ResponseWriter, r *http.Request) {
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

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Something went wrong")
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
