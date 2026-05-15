package core

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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

	cursor, limit, err := extractQueryParams(r.URL.Query())
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse query params")
		return
	}

	tldrs, err := t.Queries.GetTLDRsByUser(r.Context(), database.GetTLDRsByUserParams{
		UserID:    userId,
		CreatedAt: time.Time(cursor),
		Limit:     int64(limit),
	})
	if errors.Is(err, sql.ErrNoRows) || tldrs == nil {
		t.jsonResponse(w, http.StatusOK, Page[database.GetTLDRsByUserRow]{Results: []database.GetTLDRsByUserRow{}})
		return
	}

	if err != nil {
		t.Logger.Error("Failed to get TLDRs", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get TLDRs")
		return
	}

	if int(limit) > len(tldrs) {
		t.jsonResponse(w, http.StatusOK, Page[database.GetTLDRsByUserRow]{Results: tldrs})
		return
	}

	next := tldrs[limit-1]
	t.jsonResponse(w, http.StatusOK, Page[database.GetTLDRsByUserRow]{
		Results: tldrs[:limit-1],
		Next:    (*PageCursor)(&next.CreatedAt),
	})
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
	cursor, limit, err := extractQueryParams(r.URL.Query())
	if err != nil {
		t.errorResponse(w, r.Context(), http.StatusBadRequest, "Failed to parse query params")
		return
	}

	tldrs, err := t.Queries.GetTLDRs(r.Context(), database.GetTLDRsParams{
		CreatedAt: time.Time(cursor),
		Limit:     int64(limit),
	})
	if errors.Is(err, sql.ErrNoRows) || tldrs == nil {
		t.jsonResponse(w, http.StatusOK, Page[database.GetTLDRsRow]{Results: []database.GetTLDRsRow{}})
		return
	}

	if err != nil {
		t.Logger.Error("Failed to get TLDRs", "error", err.Error())
		t.errorResponse(w, r.Context(), http.StatusInternalServerError, "Failed to get TLDRs")
		return
	}

	if int(limit) > len(tldrs) {
		t.jsonResponse(w, http.StatusOK, Page[database.GetTLDRsRow]{Results: tldrs})
		return
	}

	next := tldrs[limit-1]
	t.jsonResponse(w, http.StatusOK, Page[database.GetTLDRsRow]{
		Results: tldrs[:limit-1],
		Next:    (*PageCursor)(&next.CreatedAt),
	})
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
	var updateRequest database.UpdateTLDRTitleByIdParams
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

	tldr, err := t.Queries.UpdateTLDRTitleById(r.Context(), updateRequest)
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

func extractQueryParams(query url.Values) (PageCursor, PageLimit, error) {
	cursorQuery := strings.TrimSpace(query.Get("cursor"))
	if cursorQuery == "" {
		cursorQuery = NoCursor
	}

	cursor, err := time.Parse(time.RFC3339, cursorQuery)
	if err != nil {
		return PageCursor{}, 0, err
	}

	limitQuery := strings.TrimSpace(query.Get("limit"))
	if limitQuery == "" {
		limitQuery = defaultPageLimit
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		return PageCursor{}, 0, err
	}
	limit += 1

	return PageCursor(cursor), PageLimit(limit), nil
}
