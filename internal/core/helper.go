package core

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/duanechan/tldr/internal/auth"
	"github.com/duanechan/tldr/internal/database"
	"github.com/google/uuid"
)

func (t *TLDR) errorResponse(w http.ResponseWriter, ctx context.Context, code int, message string, fieldErrors ...FieldError) {
	requestId, ok := ctx.Value(requestIdKey).(string)
	if !ok {
		t.Logger.Error("Failed to get request ID")
	}

	w.Header().Set("Content-Type", "application/json")

	t.jsonResponse(w, code, ErrorResponse{
		Code:      code,
		RequestID: requestId,
		Message:   message,
		Errors:    fieldErrors,
	})
}

func (t *TLDR) jsonResponse(w http.ResponseWriter, code int, payload any) {
	if code == http.StatusNoContent {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(payload)
	if err != nil {
		t.Logger.Error("Marshal error:", "error", err.Error())
		code = http.StatusInternalServerError
		body = []byte(`{"error":"Something went wrong"}`)
	}

	w.WriteHeader(code)

	if _, err := w.Write(body); err != nil {
		t.Logger.Error("Write error:", "error", err.Error())
	}
}

func (t *TLDR) insertRefreshToken(ctx context.Context, id uuid.UUID) (*database.RefreshToken, error) {
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
		UserID:    id,
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

func (t *TLDR) insertTLDR(ctx context.Context, subject, content string) (*database.Tldr, error) {
	userId, err := uuid.Parse(subject)
	if err != nil {
		return nil, err
	}

	tldrId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	tldr, err := t.Queries.CreateTLDR(ctx, database.CreateTLDRParams{
		ID:      tldrId,
		Title:   "TLDR-" + tldrId.String()[:6],
		Content: content,
		UserID:  userId,
	})
	if err != nil {
		return nil, err
	}

	return &tldr, nil
}
