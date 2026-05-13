package core

import (
	"database/sql"
	"encoding/json"
	"errors"
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

	tldr, err := t.Queries.GetTLDRById(r.Context(), database.GetTLDRByIdParams{
		UserID: userId,
		ID:     tldrId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			t.errorResponse(w, r.Context(), http.StatusNotFound, "TLDR not found")
			return
		}
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

	tldrs, err := t.Queries.GetTLDRsByUser(r.Context(), userId)
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

	updateRequest := database.UpdateTLDRTitleByIdParams{
		UserID: userId,
		ID:     tldrId,
	}
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Invalid request body")
		return
	}

	tldr, err := t.Queries.UpdateTLDRTitleById(r.Context(), updateRequest)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			t.errorResponse(w, r.Context(), http.StatusNotFound, "Failed to update TLDR with ID: "+tldrId.String())
			return
		}
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

	if err := t.Queries.DeleteTLDRById(r.Context(), database.DeleteTLDRByIdParams{
		UserID: userId,
		ID:     tldrId,
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			t.errorResponse(w, r.Context(), http.StatusNotFound, "Failed to delete TLDR with ID: "+tldrId.String())
			return
		}
		t.Logger.Error("Failed to delete TLDR", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to delete TLDR")
		return
	}

	t.jsonResponse(w, http.StatusNoContent, nil)
}
