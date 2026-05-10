package core

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (t *TLDR) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	cleanedUsername := strings.TrimSpace(req.Username)
	if cleanedUsername == "" || strings.TrimSpace(req.Password) == "" {
		errorResponse(w, http.StatusBadRequest, "Username/password is required")
		return
	}

	user, err := t.DB.GetUserByName(r.Context(), cleanedUsername)
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid username/password")
		return
	}

	matches, err := argon2id.ComparePasswordAndHash(req.Password, user.Password)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !matches {
		errorResponse(w, http.StatusUnauthorized, "Invalid username/password")
		return
	}

	jsonResponse(w, http.StatusOK, loginResponse{Token: "Sample-Token-123"})
}
