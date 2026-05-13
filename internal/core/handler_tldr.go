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
		t.errorResponse(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		t.errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	tldr, err := t.Queries.GetTLDRById(r.Context(), database.GetTLDRByIdParams{
		UserID: userId,
		ID:     tldrId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			t.errorResponse(w, http.StatusNotFound, "TLDR not found")
			return
		}
		t.errorResponse(w, http.StatusInternalServerError, "Failed to get TLDR with ID: "+tldrId.String())
		return
	}

	t.jsonResponse(w, http.StatusOK, tldr)
}

func (t *TLDR) GetTLDRs(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsKey).(*jwt.RegisteredClaims)
	if !ok {
		t.errorResponse(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		t.errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	tldrs, err := t.Queries.GetTLDRsByUser(r.Context(), userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			t.errorResponse(w, http.StatusNotFound, "TLDR not found")
			return
		}
		t.errorResponse(w, http.StatusInternalServerError, "Failed to get TLDRs")
		return
	}

	t.jsonResponse(w, http.StatusOK, tldrs)
}

func (t *TLDR) UpdateTLDR(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsKey).(*jwt.RegisteredClaims)
	if !ok {
		t.errorResponse(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		t.errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	updateRequest := database.UpdateTLDRTitleByIdParams{
		UserID: userId,
		ID:     tldrId,
	}
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		t.errorResponse(w, http.StatusInternalServerError, "Invalid request body")
		return
	}

	tldr, err := t.Queries.UpdateTLDRTitleById(r.Context(), updateRequest)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			t.errorResponse(w, http.StatusNotFound, "Failed to get TLDR with ID: "+tldrId.String())
			return
		}
		t.errorResponse(w, http.StatusInternalServerError, "Failed to update TLDR")
		return
	}

	t.jsonResponse(w, http.StatusOK, tldr)
}

func (t *TLDR) DeleteTLDR(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsKey).(*jwt.RegisteredClaims)
	if !ok {
		t.errorResponse(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		t.errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		t.errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if err := t.Queries.DeleteTLDRById(r.Context(), database.DeleteTLDRByIdParams{
		UserID: userId,
		ID:     tldrId,
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			t.errorResponse(w, http.StatusNotFound, "Failed to get TLDR with ID: "+tldrId.String())
			return
		}
		t.errorResponse(w, http.StatusInternalServerError, "Failed to delete TLDR")
		return
	}

	t.jsonResponse(w, http.StatusNoContent, nil)
}
