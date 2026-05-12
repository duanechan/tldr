package core

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/duanechan/tldr/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (t *TLDR) GetTLDR(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsKey).(*jwt.RegisteredClaims)
	if !ok {
		errorResponse(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	tldr, err := t.Queries.GetTLDRById(r.Context(), database.GetTLDRByIdParams{
		UserID: userId,
		ID:     tldrId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponse(w, http.StatusNotFound, "TLDR not found")
			return
		}
		errorResponse(w, http.StatusInternalServerError, "Failed to get TLDR with ID: "+tldrId.String())
		return
	}

	jsonResponse(w, http.StatusOK, tldr)
}

func (t *TLDR) GetTLDRs(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsKey).(*jwt.RegisteredClaims)
	if !ok {
		errorResponse(w, http.StatusUnauthorized, "Invalid claims")
		return
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	tldrs, err := t.Queries.GetTLDRsByUser(r.Context(), userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponse(w, http.StatusNotFound, "TLDR not found")
			return
		}
		errorResponse(w, http.StatusInternalServerError, "Failed to get TLDRs")
	}

	jsonResponse(w, http.StatusOK, tldrs)
}
