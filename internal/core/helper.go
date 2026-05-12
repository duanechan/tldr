package core

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/duanechan/tldr/internal/auth"
	"github.com/duanechan/tldr/internal/database"
	"github.com/google/uuid"
)

func errorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(struct {
		Error string
		Code  int
	}{
		Error: message,
		Code:  code,
	})
	if err != nil {
		log.Printf("marshal error: %v", err)
		code = http.StatusInternalServerError
		body = []byte(`{"error":"Something went wrong"}`)
	}

	w.WriteHeader(code)

	if n, err := w.Write(body); err != nil {
		log.Printf("write error: %v", err)
	} else if n < len(body) {
		log.Printf("short write: wrote %d of %d bytes", n, len(body))
	}
}

func jsonResponse(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("marshal error: %v", err)
		code = http.StatusInternalServerError
		body = []byte(`{"error":"Something went wrong"}`)
	}

	w.WriteHeader(code)

	if n, err := w.Write(body); err != nil {
		log.Printf("write error: %v", err)
	} else if n < len(body) {
		log.Printf("short write: wrote %d of %d bytes", n, len(body))
	}
}

func (t *TLDR) insertRefreshToken(ctx context.Context, user database.User) (*database.RefreshToken, error) {
	refreshTokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := auth.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshToken, err := t.Queries.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		ID:        refreshTokenId,
		Token:     refreshTokenString,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(t.Config.RefreshExpiry),
	})
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (t *TLDR) setRefreshTokenCookie(w http.ResponseWriter, refreshToken database.RefreshToken) {
	var sameSite http.SameSite

	switch t.Config.Environment {
	case "prod":
		sameSite = http.SameSiteStrictMode
	case "dev":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteLaxMode
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "REFRESH_TOKEN",
		Value:    refreshToken.Token,
		Path:     "/",
		Expires:  refreshToken.ExpiresAt,
		HttpOnly: true,
		SameSite: sameSite,
		Secure:   t.Config.Environment == "prod",
	})
}
