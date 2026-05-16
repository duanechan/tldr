package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/duanechan/tldr/internal/types"
)

func (t *TLDR) errorResponse(w http.ResponseWriter, ctx context.Context, code int, message string, fieldErrors ...types.FieldError) {
	requestId, ok := ctx.Value(types.RequestIdKey).(string)
	if !ok {
		t.Logger.Error("Failed to get request ID")
	}

	w.Header().Set("Content-Type", "application/json")

	t.jsonResponse(w, code, types.ErrorResponse{
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
