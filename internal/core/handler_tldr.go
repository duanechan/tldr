package core

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/duanechan/tldr/internal/auth"
	"github.com/duanechan/tldr/internal/database"
	"github.com/google/uuid"
)

// UserGetTLDR returns a user's TLDR from the given "id" path.
func (a *App) UserGetTLDR(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserID(r.Context())
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusUnauthorized,
			"Invalid claims",
		)
		return
	}

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to parse TLDR ID",
		)
		return
	}

	tldr, err := a.Queries.GetTLDRByIDAndUser(
		r.Context(),
		database.GetTLDRByIDAndUserParams{
			UserID: userId,
			ID:     tldrId,
		},
	)
	if errors.Is(err, sql.ErrNoRows) {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusNotFound,
			"TLDR not found",
		)
		return
	}

	if err != nil {
		a.Logger.Error("Failed to get TLDR", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to get TLDR",
		)
		return
	}

	a.jsonResponse(w, http.StatusOK, tldr)
}

// UserGetTLDRs returns a paginated response of a user's TLDRs.
func (a *App) UserGetTLDRs(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserID(r.Context())
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusUnauthorized,
			"Invalid claims",
		)
		return
	}

	createdAt, id, limit, fieldErrors := extractQueryParams(r.URL.Query())
	if fieldErrors != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to parse query params",
			fieldErrors...)
		return
	}

	tldrs, err := a.Queries.GetTLDRsByUser(
		r.Context(),
		database.GetTLDRsByUserParams{
			UserID:      userId,
			CreatedAt:   *createdAt,
			CreatedAt_2: *createdAt,
			ID:          id,
			Limit:       limit + 1,
		},
	)
	if errors.Is(err, sql.ErrNoRows) || tldrs == nil {
		a.jsonResponse(
			w,
			http.StatusOK,
			Page[database.GetTLDRsByUserRow]{
				Results: []database.GetTLDRsByUserRow{},
			},
		)
		return
	}

	if err != nil {
		a.Logger.Error("Failed to get TLDRs", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to get TLDRs",
		)
		return
	}

	if int(limit) >= len(tldrs) {
		a.jsonResponse(
			w,
			http.StatusOK,
			Page[database.GetTLDRsByUserRow]{Results: tldrs},
		)
		return
	}

	lastItem := tldrs[limit]
	next := encodeCursor(&lastItem.CreatedAt, lastItem.ID)

	a.jsonResponse(w, http.StatusOK, Page[database.GetTLDRsByUserRow]{
		Results: tldrs[:limit],
		Next:    next,
	})
}

// UserUpdateTLDR returns a TLDR with an updated title from the given
// "id" path.
func (a *App) UserUpdateTLDR(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserID(r.Context())
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusUnauthorized,
			"Failed to parse user ID",
		)
		return
	}

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to parse TLDR ID",
		)
		return
	}

	var updateRequest database.UpdateTLDRTitleParams
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	updateRequest.UserID = userId
	updateRequest.ID = tldrId

	tldr, err := a.Queries.UpdateTLDRTitle(r.Context(), updateRequest)
	if errors.Is(err, sql.ErrNoRows) {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusNotFound,
			fmt.Sprintf("TLDR with ID: %s not found", tldrId.String()),
		)
		return
	}

	if err != nil {
		a.Logger.Error("Failed to update TLDR", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to update TLDR",
		)
		return
	}

	a.jsonResponse(w, http.StatusOK, tldr)
}

// UserDeleteTLDR deletes a TLDR of a user from the given "id" path.
func (a *App) UserDeleteTLDR(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserID(r.Context())
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusUnauthorized,
			"Failed to parse user ID",
		)
		return
	}

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to parse TLDR ID",
		)
		return
	}

	if err := a.Queries.DeleteTLDR(r.Context(), database.DeleteTLDRParams{
		UserID: userId,
		ID:     tldrId,
	}); err != nil {
		a.Logger.Error("Failed to delete TLDR", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to delete TLDR",
		)
		return
	}

	a.jsonResponse(w, http.StatusNoContent, nil)
}

// UserDeleteTLDRs batch deletes TLDRs of a user from a given list of IDs
func (a *App) UserDeleteTLDRs(w http.ResponseWriter, r *http.Request) {
	userId, err := auth.GetUserID(r.Context())
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusUnauthorized,
			"Failed to parse user ID",
		)
		return
	}

	var ids []uuid.UUID
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	res, err := a.Queries.DeleteTLDRsByIdAndUser(
		r.Context(),
		database.DeleteTLDRsByIdAndUserParams{
			UserID: userId,
			Ids:    ids,
		})
	if err != nil {
		a.Logger.Error("Failed to delete TLDRs", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to delete TLDRs",
		)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		a.Logger.Error(
			"Failed to delete TLDRs",
			"error",
			err.Error(),
			"rows",
			rowsAffected,
		)
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to delete TLDRs",
		)
		return
	}

	a.jsonResponse(
		w,
		http.StatusOK,
		fmt.Sprintf("Deleted %d TLDRs", rowsAffected),
	)
}

// AdminGetTLDR returns a TLDR from the given "id" path.
func (a *App) AdminGetTLDR(w http.ResponseWriter, r *http.Request) {
	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to parse TLDR ID",
		)
		return
	}

	tldr, err := a.Queries.GetTLDRById(r.Context(), tldrId)
	if errors.Is(err, sql.ErrNoRows) {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusNotFound,
			fmt.Sprintf("TLDR with ID: %s not found", tldrId.String()),
		)
		return
	}

	if err != nil {
		a.Logger.Error("Failed to get TLDR", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to get TLDR",
		)
		return
	}

	a.jsonResponse(w, http.StatusOK, tldr)
}

// AdminGetTLDRs returns a paginated response of TLDRs.
func (a *App) AdminGetTLDRs(w http.ResponseWriter, r *http.Request) {
	createdAt, id, limit, fieldErrors := extractQueryParams(r.URL.Query())
	if fieldErrors != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to parse query params",
			fieldErrors...)
		return
	}

	tldrs, err := a.Queries.GetTLDRs(r.Context(), database.GetTLDRsParams{
		CreatedAt:   *createdAt,
		CreatedAt_2: *createdAt,
		ID:          id,
		Limit:       limit + 1,
	})
	if errors.Is(err, sql.ErrNoRows) || tldrs == nil {
		a.jsonResponse(
			w,
			http.StatusOK,
			Page[database.GetTLDRsByUserRow]{
				Results: []database.GetTLDRsByUserRow{},
			},
		)
		return
	}

	if err != nil {
		a.Logger.Error("Failed to get TLDRs", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to get TLDRs",
		)
		return
	}

	if int(limit) >= len(tldrs) {
		a.jsonResponse(
			w,
			http.StatusOK,
			Page[database.GetTLDRsRow]{Results: tldrs},
		)
		return
	}

	lastItem := tldrs[limit]
	next := encodeCursor(&lastItem.CreatedAt, lastItem.ID)

	a.jsonResponse(w, http.StatusOK, Page[database.GetTLDRsRow]{
		Results: tldrs[:limit],
		Next:    next,
	})
}

// AdminUpdateTLDR returns a TLDR with an updated title from the given
// "id" path.
func (a *App) AdminUpdateTLDR(w http.ResponseWriter, r *http.Request) {
	var updateRequest database.UpdateTLDRTitleByIdParams
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to parse TLDR ID",
		)
		return
	}

	updateRequest.ID = tldrId

	tldr, err := a.Queries.UpdateTLDRTitleById(r.Context(), updateRequest)
	if errors.Is(err, sql.ErrNoRows) {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusNotFound,
			fmt.Sprintf("TLDR with ID: %s not found", tldrId),
		)
		return
	}

	if err != nil {
		a.Logger.Error("Failed to update TLDR", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to update TLDR",
		)
		return
	}

	a.jsonResponse(w, http.StatusOK, tldr)
}

// AdminDeleteTLDR deletes a TLDR from the given "id" path.
func (a *App) AdminDeleteTLDR(w http.ResponseWriter, r *http.Request) {
	tldrId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to parse TLDR ID",
		)
		return
	}

	if err = a.Queries.DeleteTLDRById(r.Context(), tldrId); err != nil {
		a.Logger.Error("Failed to delete TLDR", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to delete TLDR",
		)
		return
	}

	a.jsonResponse(w, http.StatusNoContent, nil)
}

// AdminDeleteTLDRs batch deletes TLDRs from a given list of IDs
func (a *App) AdminDeleteTLDRs(w http.ResponseWriter, r *http.Request) {
	var ids []uuid.UUID
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	res, err := a.Queries.DeleteTLDRs(r.Context(), ids)
	if err != nil {
		a.Logger.Error("Failed to delete TLDRs", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to delete TLDRs",
		)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		a.Logger.Error(
			"Failed to delete TLDRs",
			"error",
			err.Error(),
			"rows",
			rowsAffected,
		)
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to delete TLDRs",
		)
		return
	}

	a.jsonResponse(w, http.StatusOK, struct {
		Message string
	}{
		Message: fmt.Sprintf("Deleted %d TLDRs", rowsAffected),
	})
}

// AdminDeleteAllTLDRs deletes all TLDRs from the database.
// Intended only for development.
func (a *App) AdminDeleteAllTLDRs(w http.ResponseWriter, r *http.Request) {
	if a.Config.Environment != "dev" {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusUnauthorized,
			"Invalid application environment",
		)
		return
	}

	res, err := a.Queries.DeleteAllTLDRs(r.Context())
	if err != nil {
		a.Logger.Error("Failed to delete TLDRs", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to delete TLDRs",
		)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		a.Logger.Error(
			"Failed to delete TLDRs",
			"error",
			err.Error(),
			"rows",
			rowsAffected,
		)
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to delete TLDRs",
		)
		return
	}

	a.jsonResponse(
		w,
		http.StatusOK,
		fmt.Sprintf("Deleted %d TLDRs", rowsAffected),
	)
}
