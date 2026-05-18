package core

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/duanechan/tldr/internal/auth"
	"github.com/duanechan/tldr/internal/database"
	"github.com/google/uuid"
)

func (a *App) UserGetMe(w http.ResponseWriter, r *http.Request) {
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

	a.getUser(w, r, userId)
}

func (a *App) UserUpdateUsername(w http.ResponseWriter, r *http.Request) {
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

	a.updateUsername(w, r, userId)
}

func (a *App) UserUpdatePassword(w http.ResponseWriter, r *http.Request) {
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

	a.updatePassword(w, r, userId)
}

func (a *App) AdminGetUser(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to parse user ID",
		)
		return
	}

	a.getUser(w, r, userId)
}

func (a *App) AdminGetUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := a.Queries.GetUsers(r.Context(), database.GetUsersParams{
		CreatedAt:   *createdAt,
		CreatedAt_2: *createdAt,
		ID:          id,
		Limit:       limit + 1,
	})
	if errors.Is(err, sql.ErrNoRows) || users == nil {
		a.jsonResponse(w, http.StatusOK, []database.User{})
		return
	}

	if err != nil {
		a.Logger.Error("Failed to get users", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to get users",
		)
		return
	}

	if int(limit) >= len(users) {
		a.jsonResponse(
			w,
			http.StatusOK,
			Page[database.GetUsersRow]{Results: users},
		)
		return
	}

	lastItem := users[limit]
	next := encodeCursor(&lastItem.CreatedAt, lastItem.ID)

	a.jsonResponse(w, http.StatusOK, Page[database.GetUsersRow]{
		Results: users[:limit],
		Next:    next,
	})
}

func (a *App) AdminUpdateUsername(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to parse user ID",
		)
		return
	}

	a.updateUsername(w, r, userId)
}

func (a *App) AdminUpdatePassword(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to parse user ID",
		)
		return
	}

	a.updatePassword(w, r, userId)
}

func (a *App) AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to parse user ID",
		)
		return
	}

	if err = a.Queries.DeleteUser(r.Context(), userId); err != nil {
		a.Logger.Error("Failed to delete user", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to delete user",
		)
		return
	}

	a.jsonResponse(w, http.StatusNoContent, nil)
}

func (a *App) AdminDeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	if a.Config.Environment != "dev" {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusUnauthorized,
			"Invalid application environment",
		)
		return
	}

	res, err := a.Queries.DeleteAllUsers(r.Context())
	if err != nil {
		a.Logger.Error("Failed to delete users", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to delete users",
		)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		a.Logger.Error(
			"Failed to delete users",
			"error",
			err.Error(),
			"rows",
			rowsAffected,
		)
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to delete users",
		)
		return
	}

	a.jsonResponse(w, http.StatusOK, struct {
		Message string
	}{
		Message: fmt.Sprintf("Deleted %d users", rowsAffected),
	})
}

func (a *App) getUser(
	w http.ResponseWriter,
	r *http.Request,
	userId uuid.UUID,
) {
	user, err := a.Queries.GetUserById(r.Context(), userId)
	if errors.Is(err, sql.ErrNoRows) {
		a.errorResponse(w, r.Context(), http.StatusNotFound, "User not found")
		return
	}

	if err != nil {
		a.Logger.Error("Failed to get user", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to get user",
		)
		return
	}

	a.jsonResponse(w, http.StatusOK, user)
}

func (a *App) updateUsername(
	w http.ResponseWriter,
	r *http.Request,
	userId uuid.UUID,
) {
	var updateRequest database.UpdateUsernameParams
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	updateRequest.ID = userId

	var fieldError *FieldError
	cleanedUsername := strings.TrimSpace(updateRequest.Username)
	if cleanedUsername == "" {
		fieldError = &FieldError{
			Field:   "username",
			Message: "Username is required",
		}
	} else if len(cleanedUsername) < minimumUsernameLength {
		fieldError = &FieldError{
			Field: "username",
			Message: fmt.Sprintf(
				"Username must be %d characters long",
				minimumUsernameLength,
			),
		}
	}

	if fieldError != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to update username",
			*fieldError,
		)
		return
	}

	updateRequest.Username = cleanedUsername

	user, err := a.Queries.UpdateUsername(r.Context(), updateRequest)
	if errors.Is(err, sql.ErrNoRows) {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusNotFound,
			"Failed to update user",
		)
		return
	}

	if err != nil {
		a.Logger.Error("Failed to update user", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to update user",
		)
		return
	}

	a.jsonResponse(w, http.StatusOK, user)
}

func (a *App) updatePassword(
	w http.ResponseWriter,
	r *http.Request,
	userId uuid.UUID,
) {
	var updateRequest database.UpdatePasswordParams
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	updateRequest.ID = userId

	var fieldError *FieldError
	if strings.TrimSpace(updateRequest.Password) == "" {
		fieldError = &FieldError{
			Field:   "password",
			Message: "Password is required",
		}
	} else if len(updateRequest.Password) < minimumPasswordLength {
		fieldError = &FieldError{
			Field: "password",
			Message: fmt.Sprintf(
				"Password must be %d characters long",
				minimumPasswordLength,
			),
		}
	}

	if fieldError != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusBadRequest,
			"Failed to update password",
			*fieldError,
		)
		return
	}

	hashedPassword, err := argon2id.CreateHash(
		updateRequest.Password,
		argon2id.DefaultParams,
	)
	if err != nil {
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Something went wrong",
		)
		return
	}

	updateRequest.Password = hashedPassword

	user, err := a.Queries.UpdatePassword(r.Context(), updateRequest)
	if errors.Is(err, sql.ErrNoRows) {
		a.errorResponse(w, r.Context(), http.StatusNotFound, "User not found")
		return
	}

	if err != nil {
		a.Logger.Error("Failed to update password", "error", err.Error())
		a.errorResponse(
			w,
			r.Context(),
			http.StatusInternalServerError,
			"Failed to update password",
		)
		return
	}

	a.jsonResponse(w, http.StatusOK, user)
}
